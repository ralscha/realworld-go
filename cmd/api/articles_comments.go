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
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	var commentRequest dto.CommentRequest
	if ok := request.DecodeJSONValidate[*dto.CommentRequest](w, r, &commentRequest, dto.ValidateCommentRequest); !ok {
		return
	}

	articleSlug := chi.URLParam(r, "slug")
	article, err := models.Articles(qm.Select(models.ArticleColumns.ID, models.ArticleColumns.UserID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	epochSeconds := time.Now().Unix()
	comment := models.Comment{
		Body:      commentRequest.Comment.Body,
		ArticleID: article.ID,
		UserID:    userID,
		CreatedAt: epochSeconds,
		UpdatedAt: epochSeconds,
	}

	if err := comment.Insert(r.Context(), app.db, boil.Infer()); err != nil {
		response.ServerError(w, err)
		return
	}

	user, err := models.Users(qm.Select(models.UserColumns.Username, models.UserColumns.Bio, models.UserColumns.Image), models.UserWhere.ID.EQ(article.UserID)).One(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
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
		CreatedAt: time.Unix(comment.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt: time.Unix(comment.UpdatedAt, 0).Format(time.RFC3339),
		Body:      comment.Body,
		Author:    profile,
	}
	response.JSON(w, http.StatusCreated, insertedComment)

}

func (app *application) articlesGetComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := article.Comments(qm.Load(models.CommentRels.User)).All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	followingIds, err := models.Follows(qm.Select(models.FollowColumns.FollowID), models.FollowWhere.UserID.EQ(userID)).All(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
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
			CreatedAt: time.Unix(comment.CreatedAt, 0).Format(time.RFC3339),
			UpdatedAt: time.Unix(comment.UpdatedAt, 0).Format(time.RFC3339),
			Body:      comment.Body,
			Author:    profile,
		})
	}

	response.JSON(w, http.StatusOK, dto.CommentMany{Comments: commentsResponse})
}

func (app *application) articlesDeleteComment(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)
	articleSlug := chi.URLParam(r, "slug")
	commentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.NotFound(w, r)
		return
	}

	article, err := models.Articles(qm.Select(models.ArticleColumns.ID), models.ArticleWhere.Slug.EQ(articleSlug)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	comment, err := article.Comments(models.CommentWhere.ID.EQ(int64(commentID)), models.CommentWhere.UserID.EQ(userID)).One(r.Context(), app.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	if err := comment.Delete(r.Context(), app.db); err != nil {
		response.ServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}
