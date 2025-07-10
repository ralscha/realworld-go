package main

import (
	"database/sql"
	"errors"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/go-chi/chi/v5"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
)

func (app *application) articlesFavorite(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	newFavorite := models.ArticleFavorite{
		UserID:    userID,
		ArticleID: article.ID,
	}

	err = newFavorite.Upsert(r.Context(), tx, false, []string{models.ArticleFavoriteColumns.UserID, models.ArticleFavoriteColumns.ArticleID}, boil.None(), boil.Infer())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleByID(r.Context(), article.ID, true, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{Article: updatedArticle})
}

func (app *application) articlesUnfavorite(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	articleSlug := chi.URLParam(r, "slug")

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	err = models.ArticleFavorites(models.ArticleFavoriteWhere.UserID.EQ(userID), models.ArticleFavoriteWhere.ArticleID.EQ(article.ID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	updatedArticle, err := app.getArticleByID(r.Context(), article.ID, true, userID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.ArticleOne{Article: updatedArticle})

}
