package main

import (
	"database/sql"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/request"
	"realworldgo.rasc.ch/internal/response"
)

var userNotFoundPasswordHash string

func initAuth(config config.Config) error {
	var err error
	userNotFoundPasswordHash, err = argon2id.CreateHash("userNotFoundPassword", &argon2id.Params{
		Memory:      config.Argon2.Memory,
		Iterations:  config.Argon2.Iterations,
		Parallelism: config.Argon2.Parallelism,
		SaltLength:  config.Argon2.SaltLength,
		KeyLength:   config.Argon2.KeyLength,
	})
	return err
}

func (app *application) usersLogin(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate[*dto.UserRequest](w, r, &userLoginRequest, dto.ValidateUserLoginRequest); !ok {
		return
	}

	user, err := models.Users(qm.Select(
		models.UserColumns.Email,
		models.UserColumns.Username,
		models.UserColumns.Bio,
		models.UserColumns.Image),
		models.UserWhere.Email.EQ(userLoginRequest.User.Email)).One(r.Context(), app.db)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if _, err := argon2id.ComparePasswordAndHash(userNotFoundPasswordHash, userNotFoundPasswordHash); err != nil {
				response.ServerError(w, err)
				return
			}
			response.Unauthorized(w)
		} else {
			response.ServerError(w, err)
		}
		return
	}

	match, err := argon2id.ComparePasswordAndHash(userLoginRequest.User.Password, user.Password)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	if !match {
		response.Unauthorized(w)
		return
	}

	token := ""
	var userResponse = dto.UserWrapper{
		User: dto.User{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio.String,
			Image:    user.Image.String,
		},
	}

	response.JSON(w, http.StatusOK, userResponse)
}

func (app *application) usersRegistration(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate[*dto.UserRequest](w, r, &userLoginRequest, dto.ValidateUserRegistrationRequest); !ok {
		return
	}

}
