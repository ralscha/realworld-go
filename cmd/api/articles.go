package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/request"
	"realworldgo.rasc.ch/internal/response"
	"strconv"
	"time"
)

func (app *application) articlesFeed(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("limit")

	if offsetParam == "" {
		offsetParam = "0"
	}
	if limitParam == "" {
		limitParam = "20"
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	follows, err := models.Follows(qm.Select(models.FollowColumns.FollowID), models.FollowWhere.UserID.EQ(userID)).All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	followIds := make([]int64, len(follows))
	for i, follow := range follows {
		followIds[i] = follow.FollowID
	}

	articles, err := models.Articles(models.ArticleWhere.UserID.IN(followIds), qm.Limit(limit), qm.Offset(offset)).All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	articlesCount, err := models.Articles(qm.Select(models.ArticleColumns.ID)).Count(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	articlesResponse := make([]dto.Article, len(articles))
	for i, article := range articles {
		dtoArticle, err := app.getArticle(r.Context(), article, userID)
		if err != nil {
			response.ServerError(w, err)
			return
		}

		articlesResponse[i] = dtoArticle
	}

	response.JSON(w, http.StatusOK, dto.ArticlesMany{
		Articles:      articlesResponse,
		ArticlesCount: articlesCount,
	})
}

func (app *application) articlesList(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)

	offsetParam := r.URL.Query().Get("offset")
	limitParam := r.URL.Query().Get("limit")

	if offsetParam == "" {
		offsetParam = "0"
	}
	if limitParam == "" {
		limitParam = "20"
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	tagParam := r.URL.Query().Get("tag")
	favoritedParam := r.URL.Query().Get("favorited")
	authorParam := r.URL.Query().Get("author")

	var mods []qm.QueryMod

	if tagParam != "" {
		mods = append(mods,
			qm.InnerJoin(models.TableNames.ArticleTag+" ON "+models.TableNames.ArticleTag+"."+models.ArticleTagColumns.ArticleID+" = "+models.TableNames.Article+"."+models.ArticleColumns.ID))
		mods = append(mods,
			qm.InnerJoin(models.TableNames.Tag+" ON "+models.TableNames.Tag+"."+models.TagColumns.ID+" = "+models.TableNames.ArticleTag+"."+models.ArticleTagColumns.TagID))
		mods = append(mods, models.TagWhere.Name.EQ(tagParam))
	}

	if favoritedParam != "" {
		mods = append(mods,
			qm.InnerJoin(models.TableNames.ArticleFavorite+" ON "+models.TableNames.ArticleFavorite+"."+models.ArticleFavoriteColumns.ArticleID+" = "+models.TableNames.Article+"."+models.ArticleColumns.ID))
		mods = append(mods,
			qm.InnerJoin(models.TableNames.User+" ON "+models.TableNames.User+"."+models.UserColumns.ID+" = "+models.TableNames.ArticleFavorite+"."+models.ArticleFavoriteColumns.UserID))
		mods = append(mods, models.UserWhere.Username.EQ(favoritedParam))
	}

	if authorParam != "" {
		mods = append(mods,
			qm.InnerJoin(models.TableNames.User+" ON "+models.TableNames.User+"."+models.UserColumns.ID+" = "+models.TableNames.Article+"."+models.ArticleColumns.UserID))
		mods = append(mods, models.UserWhere.Username.EQ(authorParam))
	}

	mods = append(mods, qm.Limit(limit), qm.Offset(offset))

	articlesCount, err := models.Articles(mods...).Count(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	if articlesCount == 0 {
		response.JSON(w, http.StatusOK, dto.ArticlesMany{
			Articles:      []dto.Article{},
			ArticlesCount: 0,
		})
		return
	}

	articles, err := models.Articles(mods...).All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	articlesResponse := make([]dto.Article, len(articles))
	for i, article := range articles {
		dtoArticle, err := app.getArticle(r.Context(), article, userID)
		if err != nil {
			response.ServerError(w, err)
			return
		}

		articlesResponse[i] = dtoArticle
	}

	response.JSON(w, http.StatusOK, dto.ArticlesMany{
		Articles:      articlesResponse,
		ArticlesCount: articlesCount,
	})
}

func (app *application) articleGet(w http.ResponseWriter, r *http.Request) {
	articleSlug := chi.URLParam(r, "slug")

	article, err := app.getArticleBySlug(r.Context(), articleSlug, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: article,
	})
}

func (app *application) articlesCreate(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	var articleRequest dto.ArticleRequest
	if ok := request.DecodeJSONValidate[*dto.ArticleRequest](w, r, &articleRequest, dto.ValidateArticleCreateRequest); !ok {
		return
	}

	secondsSinceEpoch := time.Now().Unix()
	newArticle := models.Article{
		Title:       articleRequest.Article.Title,
		Description: articleRequest.Article.Description,
		Body:        articleRequest.Article.Body,
		UserID:      userID,
		Slug:        slug.Make(articleRequest.Article.Title),
		CreatedAt:   secondsSinceEpoch,
		UpdatedAt:   secondsSinceEpoch,
	}

	err := newArticle.Insert(r.Context(), app.db, boil.Infer())
	if err != nil {
		response.ServerError(w, err)
		return
	}

	if len(articleRequest.Article.TagList) > 0 {
		tags, err := models.Tags(models.TagWhere.Name.IN(articleRequest.Article.TagList)).All(r.Context(), app.db)
		if err != nil {
			response.ServerError(w, err)
			return
		}

		for _, tag := range tags {
			articleTag := models.ArticleTag{
				ArticleID: newArticle.ID,
				TagID:     tag.ID,
			}
			err = articleTag.Insert(r.Context(), app.db, boil.Infer())
			if err != nil {
				response.ServerError(w, err)
				return
			}
		}
	}

	insertedArticle, err := app.getArticleById(r.Context(), newArticle.ID, userID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: insertedArticle,
	})
}

func (app *application) articlesUpdate(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	articleSlug := chi.URLParam(r, "slug")

	var articleRequest dto.ArticleRequest
	if err := request.DecodeJSON(w, r, articleRequest); err != nil {
		response.BadRequest(w, err)
		return
	}

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug), models.ArticleWhere.UserID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	updates := models.M{
		models.ArticleColumns.UpdatedAt: time.Now().Unix(),
	}
	if articleRequest.Article.Title != "" {
		updates[models.ArticleColumns.Title] = articleRequest.Article.Title
		updates[models.ArticleColumns.Slug] = slug.Make(articleRequest.Article.Title)
	}
	if articleRequest.Article.Description != "" {
		updates[models.ArticleColumns.Description] = articleRequest.Article.Description
	}
	if articleRequest.Article.Body != "" {
		updates[models.ArticleColumns.Body] = articleRequest.Article.Body
	}

	err = models.Articles(models.ArticleWhere.ID.EQ(article.ID)).UpdateAll(r.Context(), app.db, updates)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleById(r.Context(), article.ID, userID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: updatedArticle,
	})
}

func (app *application) articlesDelete(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug), models.ArticleWhere.UserID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	err = models.Articles(models.ArticleWhere.ID.EQ(article.ID)).DeleteAll(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (app *application) getArticleById(ctx context.Context, articleId, userID int64) (dto.Article, error) {
	article, err := models.Articles(
		qm.Select(
			models.ArticleColumns.ID,
			models.ArticleColumns.UserID,
			models.ArticleColumns.Title,
			models.ArticleColumns.Description,
			models.ArticleColumns.Body,
			models.ArticleColumns.Slug,
			models.ArticleColumns.CreatedAt,
			models.ArticleColumns.UpdatedAt,
		),
		models.ArticleWhere.ID.EQ(articleId)).One(ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(ctx, article, userID)
}

func (app *application) getArticleBySlug(ctx context.Context, articleSlug string, userID int64) (dto.Article, error) {
	article, err := models.Articles(
		qm.Select(
			models.ArticleColumns.ID,
			models.ArticleColumns.UserID,
			models.ArticleColumns.Title,
			models.ArticleColumns.Description,
			models.ArticleColumns.Body,
			models.ArticleColumns.Slug,
			models.ArticleColumns.CreatedAt,
			models.ArticleColumns.UpdatedAt,
		), models.ArticleWhere.Slug.EQ(articleSlug)).One(ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(ctx, article, userID)
}

func (app *application) getArticle(ctx context.Context, article *models.Article, userID int64) (dto.Article, error) {
	author, err := models.Users(qm.Select(models.UserColumns.Username,
		models.UserColumns.Bio, models.UserColumns.Image),
		models.UserWhere.ID.EQ(article.UserID)).One(ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	following := false
	if userID != 0 {
		following, err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(author.ID)).
			Exists(ctx, app.db)
		if err != nil {
			return dto.Article{}, err
		}
	}

	authorProfile := dto.Profile{
		Username:  author.Username,
		Bio:       author.Bio.String,
		Image:     author.Image.String,
		Following: following,
	}

	favorited := false
	if userID != 0 {
		favorited, err = models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userID),
			models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
			Exists(ctx, app.db)
		if err != nil {
			return dto.Article{}, err
		}
	}

	tags, err := models.Tags(qm.Select(models.TagColumns.Name),
		qm.InnerJoin(models.TableNames.ArticleTag+" ON "+models.TableNames.ArticleTag+"."+models.ArticleTagColumns.TagID+" = "+models.TableNames.Tag+"."+models.TagColumns.ID),
		models.ArticleTagWhere.ArticleID.EQ(article.ID)).
		All(ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	tagList := make([]string, len(tags))
	for _, tag := range tags {
		tagList = append(tagList, tag.Name)
	}

	favoritesCount, err := models.ArticleFavorites(models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
		Count(ctx, app.db)
	if err != nil {
		return dto.Article{}, err
	}

	createdAtTime := time.Unix(article.CreatedAt, 0)
	updatedAtTime := time.Unix(article.UpdatedAt, 0)
	articleDto := dto.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      createdAtTime.Format(time.RFC3339),
		UpdatedAt:      updatedAtTime.Format(time.RFC3339),
		Favorited:      favorited,
		FavoritesCount: int(favoritesCount),
		Author:         authorProfile,
	}

	return articleDto, nil
}
