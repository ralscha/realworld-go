package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/gobuffalo/validate"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
		models.UserColumns.Password,
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

	token, done := app.createToken(w, r, err, user.ID)
	if done {
		return
	}

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

func (app *application) createToken(w http.ResponseWriter, r *http.Request, err error, userId int64) (string, bool) {
	ctx, err := app.sessionManager.Load(r.Context(), "")
	if err != nil {
		response.ServerError(w, err)
		return "", true
	}
	app.sessionManager.Put(ctx, "userID", userId)
	token, _, err := app.sessionManager.Commit(ctx)
	if err != nil {
		response.ServerError(w, err)
		return "", true
	}
	return token, false
}

func (app *application) usersRegistration(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate[*dto.UserRequest](w, r, &userLoginRequest, dto.ValidateUserRegistrationRequest); !ok {
		return
	}

	usernameExists, err := models.Users(models.UserWhere.Username.EQ(userLoginRequest.User.Username)).Exists(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}
	if usernameExists {
		validationError := validate.Errors{
			Errors: map[string][]string{"username": {"exists"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	emailExists, err := models.Users(models.UserWhere.Email.EQ(userLoginRequest.User.Email)).Exists(r.Context(), app.db)
	if err != nil {
		response.ServerError(w, err)
		return
	}
	if emailExists {
		validationError := validate.Errors{
			Errors: map[string][]string{"email": {"exists"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	hashedPassword, err := argon2id.CreateHash(userLoginRequest.User.Password, &argon2id.Params{
		Memory:      app.config.Argon2.Memory,
		Iterations:  app.config.Argon2.Iterations,
		Parallelism: app.config.Argon2.Parallelism,
		SaltLength:  app.config.Argon2.SaltLength,
		KeyLength:   app.config.Argon2.KeyLength,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	newUser := models.User{
		Username: userLoginRequest.User.Username,
		Password: hashedPassword,
		Email:    userLoginRequest.User.Email,
	}

	err = newUser.Insert(r.Context(), app.db, boil.Infer())
	if err != nil {
		response.ServerError(w, err)
		return
	}

	fmt.Println("userID", newUser.ID)
	token, done := app.createToken(w, r, err, newUser.ID)
	if done {
		return
	}

	var userResponse = dto.UserWrapper{
		User: dto.User{
			Email:    newUser.Email,
			Token:    token,
			Username: newUser.Username,
			Bio:      newUser.Bio.String,
			Image:    newUser.Image.String,
		},
	}

	response.JSON(w, http.StatusCreated, userResponse)

}
