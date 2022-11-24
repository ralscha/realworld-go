package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) articlesFavorite(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	newFavorite := models.ArticleFavorite{
		UserID:    userID,
		ArticleID: article.ID,
	}

	err = newFavorite.Upsert(r.Context(), app.db, false, []string{models.ArticleFavoriteColumns.UserID, models.ArticleFavoriteColumns.ArticleID}, boil.None(), boil.Infer())
	if err != nil {
		response.ServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleByID(r.Context(), article.ID, true, userID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{Article: updatedArticle})
}

func (app *application) articlesUnfavorite(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	err = models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userID), models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).DeleteAll(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleByID(r.Context(), article.ID, true, userID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{Article: updatedArticle})

}
