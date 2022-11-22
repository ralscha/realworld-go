package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
	"time"
)

func (app *application) articlesFavorite(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	slug := chi.URLParam(r, "slug")

	/*
			var record = this.dsl.select(ARTICLE.ID).from(ARTICLE)
				.where(ARTICLE.SLUG.eq(slug)).fetchOne();
		if (record != null) {
			var newRecord = this.dsl.newRecord(ARTICLE_FAVORITE);
			newRecord.setArticleId(record.get(ARTICLE.ID));
			newRecord.setUserId(user.getId());
			newRecord.store();
		}

		return ResponseEntity.ok()
				.body(Map.of("article", Util.getArticle(this.dsl, slug, user.getId())));
	*/

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if _, err := article.Favorite(r.Context(), app.db, user.ID); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articlesUnfavorite(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(slug),
		models.ArticleWhere.UserID.EQ(null.NewInt64(userID, true))).
		One(r.Context(), app.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	err = models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userID),
		models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).DeleteAll(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}
	/*


		return ResponseEntity.ok()
				.body(Map.of("article", Util.getArticle(this.dsl, slug, user.getId())));
	*/

	var articleResponse = dto.ArticleOne{
		Article: dto.Article{
			Slug:           "",
			Title:          "",
			Description:    "",
			Body:           "",
			TagList:        nil,
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
			Favorited:      false,
			FavoritesCount: 0,
			Author:         dto.Author{},
		},
	}

	response.JSON(w, http.StatusOK, articleResponse)
}
