// Package main contains the article handlers for the RealWorld API.
package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
		response.BadRequest(w, errors.New("offset must be a non-negative integer"))
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.BadRequest(w, errors.New("limit must be a non-negative integer"))
		return
	}

	if offset < 0 || limit < 0 {
		response.BadRequest(w, errors.New("offset and limit must be non-negative"))
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
	if len(followIDs) == 0 {
		response.JSON(w, http.StatusOK, dto.ArticlesMany{
			Articles:      []dto.Article{},
			ArticlesCount: 0,
		})
		return
	}

	articlesCount, err := models.Articles(models.ArticleWhere.UserID.IN(followIDs)).Count(r.Context(), tx)
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

	articles, err := models.Articles(
		models.ArticleWhere.UserID.IN(followIDs),
		qm.OrderBy(models.ArticleColumns.CreatedAt+" DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	articlesResponse, err := app.getArticles(r.Context(), articles, true, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
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
		response.BadRequest(w, errors.New("offset must be a non-negative integer"))
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		response.BadRequest(w, errors.New("limit must be a non-negative integer"))
		return
	}

	if offset < 0 || limit < 0 {
		response.BadRequest(w, errors.New("offset and limit must be non-negative"))
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

	queryMods := append([]qm.QueryMod{}, mods...)
	queryMods = append(queryMods, qm.OrderBy(models.ArticleColumns.CreatedAt+" DESC"), qm.Limit(limit), qm.Offset(offset))

	articles, err := models.Articles(queryMods...).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	articlesResponse, err := app.getArticles(r.Context(), articles, authentiated, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
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
	if ok := request.DecodeJSONValidate(w, r, &articleRequest, dto.ValidateArticleCreateRequest); !ok {
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
	articles, err := app.getArticles(ctx, models.ArticleSlice{article}, authenticated, userID)
	if err != nil {
		return dto.Article{}, err
	}
	if len(articles) == 0 {
		return dto.Article{}, sql.ErrNoRows
	}

	return articles[0], nil
}

func (app *application) getArticles(ctx context.Context, articles models.ArticleSlice, authenticated bool, userID int64) ([]dto.Article, error) {
	tx := ctx.Value(transactionKey).(*sql.Tx)
	if len(articles) == 0 {
		return []dto.Article{}, nil
	}

	articleIDs := make([]int64, 0, len(articles))
	authorIDs := make([]int64, 0, len(articles))
	seenAuthors := make(map[int64]struct{}, len(articles))
	for _, article := range articles {
		articleIDs = append(articleIDs, article.ID)
		if _, ok := seenAuthors[article.UserID]; !ok {
			seenAuthors[article.UserID] = struct{}{}
			authorIDs = append(authorIDs, article.UserID)
		}
	}

	authors, err := models.AppUsers(
		qm.Select(
			models.AppUserColumns.ID,
			models.AppUserColumns.Username,
			models.AppUserColumns.Bio,
			models.AppUserColumns.Image,
		),
		models.AppUserWhere.ID.IN(authorIDs),
	).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	authorsByID := make(map[int64]*models.AppUser, len(authors))
	for _, author := range authors {
		authorsByID[author.ID] = author
	}

	followingByAuthorID := make(map[int64]bool, len(authorIDs))
	if userID != 0 && len(authorIDs) > 0 {
		follows, err := models.Follows(
			qm.Select(models.FollowColumns.FollowID),
			models.FollowWhere.UserID.EQ(userID),
			models.FollowWhere.FollowID.IN(authorIDs),
		).All(ctx, tx)
		if err != nil {
			return nil, err
		}
		for _, follow := range follows {
			followingByAuthorID[follow.FollowID] = true
		}
	}

	favoritesCountByArticleID, err := favoriteCountsByArticleID(ctx, tx, articleIDs)
	if err != nil {
		return nil, err
	}

	favoritedByArticleID := make(map[int64]bool, len(articleIDs))
	if authenticated {
		favoritedByArticleID, err = favoritedArticleIDs(ctx, tx, userID, articleIDs)
		if err != nil {
			return nil, err
		}
	}

	articleTags, err := models.ArticleTags(
		qm.Select(models.ArticleTagColumns.ArticleID, models.ArticleTagColumns.TagID),
		qm.Load(models.ArticleTagRels.Tag, qm.Select(models.TagColumns.ID, models.TagColumns.Name)),
		models.ArticleTagWhere.ArticleID.IN(articleIDs),
	).All(ctx, tx)
	if err != nil {
		return nil, err
	}

	tagsByArticleID := make(map[int64][]string, len(articleIDs))
	for _, articleTag := range articleTags {
		if articleTag.R == nil || articleTag.R.Tag == nil {
			continue
		}
		tagsByArticleID[articleTag.ArticleID] = append(tagsByArticleID[articleTag.ArticleID], articleTag.R.Tag.Name)
	}
	for _, tagList := range tagsByArticleID {
		sort.Strings(tagList)
	}

	articlesResponse := make([]dto.Article, len(articles))
	for i, article := range articles {
		author, ok := authorsByID[article.UserID]
		if !ok {
			return nil, sql.ErrNoRows
		}

		articlesResponse[i] = dto.Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagsByArticleID[article.ID],
			CreatedAt:      article.CreatedAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt:      article.UpdatedAt.UTC().Format(time.RFC3339Nano),
			Favorited:      favoritedByArticleID[article.ID],
			FavoritesCount: favoritesCountByArticleID[article.ID],
			Author: dto.Profile{
				Username:  author.Username,
				Bio:       author.Bio.String,
				Image:     author.Image.String,
				Following: followingByAuthorID[article.UserID],
			},
		}
		if articlesResponse[i].TagList == nil {
			articlesResponse[i].TagList = []string{}
		}
	}

	return articlesResponse, nil
}

func favoriteCountsByArticleID(ctx context.Context, tx *sql.Tx, articleIDs []int64) (map[int64]int, error) {
	favoritesCountByArticleID := make(map[int64]int, len(articleIDs))
	if len(articleIDs) == 0 {
		return favoritesCountByArticleID, nil
	}

	query := "SELECT article_id, count(*) FROM article_favorite WHERE article_id IN (" +
		postgresPlaceholders(1, len(articleIDs)) +
		") GROUP BY article_id"
	rows, err := tx.QueryContext(ctx, query, int64Args(articleIDs)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var articleID int64
		var count int
		if err := rows.Scan(&articleID, &count); err != nil {
			return nil, err
		}
		favoritesCountByArticleID[articleID] = count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoritesCountByArticleID, nil
}

func favoritedArticleIDs(ctx context.Context, tx *sql.Tx, userID int64, articleIDs []int64) (map[int64]bool, error) {
	favoritedByArticleID := make(map[int64]bool, len(articleIDs))
	if len(articleIDs) == 0 {
		return favoritedByArticleID, nil
	}

	args := make([]any, 0, len(articleIDs)+1)
	args = append(args, userID)
	args = append(args, int64Args(articleIDs)...)

	query := "SELECT article_id FROM article_favorite WHERE user_id = $1 AND article_id IN (" +
		postgresPlaceholders(2, len(articleIDs)) +
		")"
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var articleID int64
		if err := rows.Scan(&articleID); err != nil {
			return nil, err
		}
		favoritedByArticleID[articleID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoritedByArticleID, nil
}

func postgresPlaceholders(start, count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(start+i)
	}
	return strings.Join(placeholders, ",")
}

func int64Args(values []int64) []any {
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return args
}
