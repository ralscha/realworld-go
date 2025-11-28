// Package main contains the article handlers for the RealWorld API.
package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/request"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) articlesFeed(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
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
		response.InternalServerError(w, err)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	follows, err := models.Follows(qm.Select(models.FollowColumns.FollowID), models.FollowWhere.UserID.EQ(userID)).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	followIDs := make([]int64, len(follows))
	for i, follow := range follows {
		followIDs[i] = follow.FollowID
	}

	articles, err := models.Articles(models.ArticleWhere.UserID.IN(followIDs), qm.Limit(limit), qm.Offset(offset)).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	articlesCount, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.UserID.IN(followIDs)).Count(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	articlesResponse := make([]dto.Article, len(articles))
	for i, article := range articles {
		dtoArticle, err := app.getArticle(r.Context(), article, true, userID)
		if err != nil {
			response.InternalServerError(w, err)
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
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	authentiated := false
	var userID int64
	if app.sessionManager.Exists(r.Context(), "userID") {
		authentiated = true
		userID = app.sessionManager.GetInt64(r.Context(), "userID")
	}

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
		response.InternalServerError(w, err)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.InternalServerError(w, err)
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
			qm.InnerJoin(models.TableNames.AppUser+" ON "+models.TableNames.AppUser+"."+models.AppUserColumns.ID+" = "+models.TableNames.ArticleFavorite+"."+models.ArticleFavoriteColumns.UserID))
		mods = append(mods, models.AppUserWhere.Username.EQ(favoritedParam))
	}

	if authorParam != "" {
		mods = append(mods,
			qm.InnerJoin(models.TableNames.AppUser+" ON "+models.TableNames.AppUser+"."+models.AppUserColumns.ID+" = "+models.TableNames.Article+"."+models.ArticleColumns.UserID))
		mods = append(mods, models.AppUserWhere.Username.EQ(authorParam))
	}

	mods = append(mods, qm.Limit(limit), qm.Offset(offset))

	articlesCount, err := models.Articles(mods...).Count(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if articlesCount == 0 {
		response.JSON(w, http.StatusOK, dto.ArticlesMany{
			Articles:      []dto.Article{},
			ArticlesCount: 0,
		})
		return
	}

	articles, err := models.Articles(mods...).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	articlesResponse := make([]dto.Article, len(articles))
	for i, article := range articles {
		dtoArticle, err := app.getArticle(r.Context(), article, authentiated, userID)
		if err != nil {
			response.InternalServerError(w, err)
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
	authentiated := false
	var userID int64
	if app.sessionManager.Exists(r.Context(), "userID") {
		authentiated = true
		userID = app.sessionManager.GetInt64(r.Context(), "userID")
	}

	articleSlug := chi.URLParam(r, "slug")

	article, err := app.getArticleBySlug(r.Context(), articleSlug, authentiated, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: article,
	})
}

func (app *application) articlesCreate(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	var articleRequest dto.ArticleRequest
	if ok := request.DecodeJSONValidate[*dto.ArticleRequest](w, r, &articleRequest, dto.ValidateArticleCreateRequest); !ok {
		return
	}

	now := time.Now()
	newArticle := models.Article{
		Title:       articleRequest.Article.Title,
		Description: articleRequest.Article.Description,
		Body:        articleRequest.Article.Body,
		UserID:      userID,
		Slug:        slug.Make(articleRequest.Article.Title),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := newArticle.Insert(r.Context(), tx, boil.Infer())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	for _, articleTag := range articleRequest.Article.TagList {
		if articleTag == "" {
			continue
		}
		tag, err := models.Tags(models.TagWhere.Name.EQ(articleTag)).One(r.Context(), tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				tag = &models.Tag{
					Name: articleTag,
				}
				err = tag.Insert(r.Context(), tx, boil.Infer())
				if err != nil {
					response.InternalServerError(w, err)
					return
				}
			} else {
				response.InternalServerError(w, err)
				return
			}
		}
		articleTag := models.ArticleTag{
			ArticleID: newArticle.ID,
			TagID:     tag.ID,
		}
		err = articleTag.Insert(r.Context(), tx, boil.Infer())
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
	}

	insertedArticle, err := app.getArticleByID(r.Context(), newArticle.ID, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: insertedArticle,
	})
}

func (app *application) articlesUpdate(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	articleSlug := chi.URLParam(r, "slug")

	var articleRequest dto.ArticleRequest
	if err := request.DecodeJSON(w, r, &articleRequest); err != nil {
		response.BadRequest(w, err)
		return
	}

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug), models.ArticleWhere.UserID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	updates := models.M{
		models.ArticleColumns.UpdatedAt: time.Now(),
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

	err = models.Articles(models.ArticleWhere.ID.EQ(article.ID)).UpdateAll(r.Context(), tx, updates)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleByID(r.Context(), article.ID, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{
		Article: updatedArticle,
	})
}

func (app *application) articlesDelete(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug), models.ArticleWhere.UserID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	err = models.Articles(models.ArticleWhere.ID.EQ(article.ID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (app *application) getArticleByID(ctx context.Context, articleID int64, userID int64) (dto.Article, error) {
	tx := ctx.Value(transactionKey).(*sql.Tx)
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
		models.ArticleWhere.ID.EQ(articleID)).One(ctx, tx)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(ctx, article, true, userID)
}

func (app *application) getArticleBySlug(ctx context.Context, articleSlug string, authenticated bool, userID int64) (dto.Article, error) {
	tx := ctx.Value(transactionKey).(*sql.Tx)
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
		), models.ArticleWhere.Slug.EQ(articleSlug)).One(ctx, tx)
	if err != nil {
		return dto.Article{}, err
	}

	return app.getArticle(ctx, article, authenticated, userID)
}

func (app *application) getArticle(ctx context.Context, article *models.Article, authenticated bool, userID int64) (dto.Article, error) {
	tx := ctx.Value(transactionKey).(*sql.Tx)
	author, err := models.AppUsers(qm.Select(models.AppUserColumns.Username,
		models.AppUserColumns.Bio, models.AppUserColumns.Image),
		models.AppUserWhere.ID.EQ(article.UserID)).One(ctx, tx)
	if err != nil {
		return dto.Article{}, err
	}

	following := false
	if userID != 0 {
		following, err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(author.ID)).
			Exists(ctx, tx)
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
	if authenticated {
		favorited, err = models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userID),
			models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
			Exists(ctx, tx)
		if err != nil {
			return dto.Article{}, err
		}
	}

	tags, err := models.Tags(qm.Select(models.TagColumns.Name),
		qm.InnerJoin(models.TableNames.ArticleTag+" ON "+models.TableNames.ArticleTag+"."+models.ArticleTagColumns.TagID+" = "+models.TableNames.Tag+"."+models.TagColumns.ID),
		models.ArticleTagWhere.ArticleID.EQ(article.ID), qm.OrderBy(models.TagColumns.Name)).
		All(ctx, tx)
	if err != nil {
		return dto.Article{}, err
	}

	tagList := make([]string, len(tags))
	for ix, tag := range tags {
		tagList[ix] = tag.Name
	}

	favoritesCount, err := models.ArticleFavorites(models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).
		Count(ctx, tx)
	if err != nil {
		return dto.Article{}, err
	}

	articleDto := dto.Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:      article.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Favorited:      favorited,
		FavoritesCount: int(favoritesCount),
		Author:         authorProfile,
	}

	return articleDto, nil
}
