package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/response"
	"strconv"
)

func (app *application) articlesFeed(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	articles, err := models.Articles().Feed(r.Context(), app.db, user.ID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, articles)
}

func (app *application) articlesGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articleGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articlesCreate(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	var input models.ArticleCreateInput
	if err := app.readJSON(w, r, &input); err != nil {
		return
	}

	article, err := input.Create(r.Context(), app.db, user.ID)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, article)
}

func (app *application) articlesUpdate(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	var input models.ArticleUpdateInput
	if err := app.readJSON(w, r, &input); err != nil {
		return
	}

	article, err := input.Update(r.Context(), app.db, user.ID, slug)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, article)
}

func (app *application) articlesDelete(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if article.AuthorID != user.ID {
		response.NotFound(w, r)
		return
	}

	if _, err := article.Delete(r.Context(), app.db); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

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

func (app *application) articlesFavorite(w http.ResponseWriter, r *http.Request) {
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

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
	user := app.sessionManager.Get(r.Context(), "user").(models.User)

	slug := chi.URLParam(r, "slug")

	article, err := models.Articles(models.ArticleWhere.Slug.EQ(slug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if _, err := article.Unfavorite(r.Context(), app.db, user.ID); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, article)
}
