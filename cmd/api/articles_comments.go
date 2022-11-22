package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
	"strconv"
)

func (app *application) articlesAddComment(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	var input models.CommentCreateInput
	if err := app.readJSON(w, r, &input); err != nil {
		return
	}

	comment, err := input.Create(r.Context(), app.db, user.ID, slug)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, comment)
}

func (app *application) articlesGetComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	comments, err := models.Comments().Get(r.Context(), app.db, slug)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, comments)
}

func (app *application) articlesDeleteComment(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.NotFound(w, r)
		return
	}

	comment, err := models.Comments(models.CommentWhere.ID.EQ(id), models.CommentWhere.Article.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if comment.AuthorID != user.ID {
		response.NotFound(w, r)
		return
	}

	if _, err := comment.Delete(r.Context(), app.db); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
