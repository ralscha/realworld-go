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

func (app *application) profilesGet(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	authentiated := false
	var userID int64
	if app.sessionManager.Exists(r.Context(), "userID") {
		authentiated = true
		userID = app.sessionManager.GetInt64(r.Context(), "userID")
	}

	username := chi.URLParam(r, "username")

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.Username.EQ(username)).One(r.Context(), tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	following := false

	if authentiated {
		following, err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(user.ID)).
			Exists(r.Context(), tx)

		if err != nil {
			response.InternalServerError(w, err)
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
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	username := chi.URLParam(r, "username")
	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.Username.EQ(username)).One(r.Context(), tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	newFollowing := models.Follow{
		UserID:   userID,
		FollowID: user.ID,
	}

	err = newFollowing.Upsert(r.Context(), tx, false, []string{models.FollowColumns.UserID, models.FollowColumns.FollowID},
		boil.None(), boil.Infer())
	if err != nil {
		response.InternalServerError(w, err)
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
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	username := chi.URLParam(r, "username")
	userID := app.sessionManager.GetInt64(r.Context(), "userID")

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.Username.EQ(username)).One(r.Context(), tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	err = models.Follows(models.FollowWhere.UserID.EQ(userID), models.FollowWhere.FollowID.EQ(user.ID)).DeleteAll(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
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
