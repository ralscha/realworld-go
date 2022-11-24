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

func (app *application) profilesGet(w http.ResponseWriter, r *http.Request) {
	authentiated := false
	var userID int64
	if app.sessionManager.Exists(r.Context(), "userID") {
		authentiated = true
		userID = app.sessionManager.GetInt64(r.Context(), "userID")
	}

	username := chi.URLParam(r, "username")

	user, err := models.Users(qm.Select(
		models.UserColumns.ID,
		models.UserColumns.Username,
		models.UserColumns.Bio,
		models.UserColumns.Image),
		models.UserWhere.Username.EQ(username)).One(r.Context(), app.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	following := false

	if authentiated {
		following, err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(user.ID)).
			Exists(r.Context(), app.db)

		if err != nil {
			response.ServerError(w, err)
			return
		}
	}

	profile := dto.ProfileWrapper{
		Profile: dto.Profile{
			Username:  user.Username,
			Bio:       user.Bio.String,
			Image:     user.Image.String,
			Following: following,
		},
	}

	response.JSON(w, http.StatusOK, profile)
}

func (app *application) profilesFollow(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)

	user, err := models.Users(qm.Select(
		models.UserColumns.ID,
		models.UserColumns.Username,
		models.UserColumns.Bio,
		models.UserColumns.Image),
		models.UserWhere.Username.EQ(username)).One(r.Context(), app.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	newFollowing := models.Follow{
		UserID:   userID,
		FollowID: user.ID,
	}

	err = newFollowing.Upsert(r.Context(), app.db, false, []string{models.FollowColumns.UserID, models.FollowColumns.FollowID},
		boil.None(), boil.Infer())
	if err != nil {
		response.ServerError(w, err)
		return
	}

	profile := dto.ProfileWrapper{
		Profile: dto.Profile{
			Username:  user.Username,
			Bio:       user.Bio.String,
			Image:     user.Image.String,
			Following: true,
		},
	}

	response.JSON(w, http.StatusOK, profile)
}

func (app *application) profilesUnfollow(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	userID := app.sessionManager.Get(r.Context(), "userID").(int64)

	user, err := models.Users(qm.Select(
		models.UserColumns.ID,
		models.UserColumns.Username,
		models.UserColumns.Bio,
		models.UserColumns.Image),
		models.UserWhere.Username.EQ(username)).One(r.Context(), app.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(user.ID)).DeleteAll(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	profile := dto.ProfileWrapper{
		Profile: dto.Profile{
			Username:  user.Username,
			Bio:       user.Bio.String,
			Image:     user.Image.String,
			Following: false,
		},
	}

	response.JSON(w, http.StatusOK, profile)
}
