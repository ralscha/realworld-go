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
	"realworldgo.rasc.ch/internal/request"
	"realworldgo.rasc.ch/internal/response"
	"strconv"
	"time"
)

func (app *application) articlesAddComment(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value("tx").(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	var commentRequest dto.CommentRequest
	if ok := request.DecodeJSONValidate[*dto.CommentRequest](w, r, &commentRequest, dto.ValidateCommentRequest); !ok {
		return
	}

	articleSlug := chi.URLParam(r, "slug")
	article, err := models.Articles(qm.Select(models.ArticleColumns.ID, models.ArticleColumns.UserID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	now := time.Now()
	comment := models.Comment{
		Body:      commentRequest.Comment.Body,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := comment.Insert(r.Context(), tx, boil.Infer()); err != nil {
		response.InternalServerError(w, err)
		return
	}

	user, err := models.AppUsers(qm.Select(models.AppUserColumns.Username, models.AppUserColumns.Bio, models.AppUserColumns.Image), models.AppUserWhere.ID.EQ(article.UserID)).One(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	profile := dto.Profile{
		Username:  user.Username,
		Bio:       user.Bio.String,
		Image:     user.Image.String,
		Following: false,
	}
	insertedComment := dto.Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: comment.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Body:      comment.Body,
		Author:    profile,
	}
	response.JSON(w, http.StatusCreated, dto.CommentOne{Comment: insertedComment})

}

func (app *application) articlesGetComments(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value("tx").(*sql.Tx)
	authentiated := false
	var userID int64
	if app.sessionManager.Exists(r.Context(), "userID") {
		authentiated = true
		userID = app.sessionManager.GetInt64(r.Context(), "userID")
	}

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

	comments, err := article.Comments(qm.Load(models.CommentRels.User)).All(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	var followingIds models.FollowSlice
	if authentiated {
		followingIds, err = models.Follows(qm.Select(models.FollowColumns.FollowID), models.FollowWhere.UserID.EQ(userID)).All(r.Context(), tx)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
	}

	followingMap := make(map[int64]bool)
	for _, following := range followingIds {
		followingMap[following.FollowID] = true
	}

	var commentsResponse []dto.Comment
	for _, comment := range comments {
		profile := dto.Profile{
			Username:  comment.R.User.Username,
			Bio:       comment.R.User.Bio.String,
			Image:     comment.R.User.Image.String,
			Following: followingMap[comment.R.User.ID],
		}
		commentsResponse = append(commentsResponse, dto.Comment{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt.UTC().Format(time.RFC3339Nano),
			UpdatedAt: comment.UpdatedAt.UTC().Format(time.RFC3339Nano),
			Body:      comment.Body,
			Author:    profile,
		})
	}

	response.JSON(w, http.StatusOK, dto.CommentMany{Comments: commentsResponse})
}

func (app *application) articlesDeleteComment(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value("tx").(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	articleSlug := chi.URLParam(r, "slug")
	commentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.NotFound(w, r)
		return
	}

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	comment, err := article.Comments(models.CommentWhere.ID.EQ(int64(commentID)), models.CommentWhere.UserID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	if err := comment.Delete(r.Context(), tx); err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
