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
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate[*dto.UserRequest](w, r, &userLoginRequest, dto.ValidateUserLoginRequest); !ok {
		return
	}

	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.ID,
		models.AppUserColumns.Email,
		models.AppUserColumns.Password,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.Email.EQ(userLoginRequest.User.Email)).One(r.Context(), tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if _, err := argon2id.ComparePasswordAndHash(userNotFoundPasswordHash, userNotFoundPasswordHash); err != nil {
				response.InternalServerError(w, err)
				return
			}
			response.Unauthorized(w)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	match, err := argon2id.ComparePasswordAndHash(userLoginRequest.User.Password, user.Password)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if !match {
		response.Unauthorized(w)
		return
	}

	token, done := app.createToken(w, r, user.ID)
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

func (app *application) createToken(w http.ResponseWriter, r *http.Request, userID int64) (string, bool) {
	ctx, err := app.sessionManager.Load(r.Context(), "")
	if err != nil {
		response.InternalServerError(w, err)
		return "", true
	}
	app.sessionManager.Put(ctx, "userID", userID)
	token, _, err := app.sessionManager.Commit(ctx)
	if err != nil {
		response.InternalServerError(w, err)
		return "", true
	}
	return token, false
}

func (app *application) usersRegistration(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate[*dto.UserRequest](w, r, &userLoginRequest, dto.ValidateUserRegistrationRequest); !ok {
		return
	}

	usernameExists, err := models.AppUsers(models.AppUserWhere.Username.EQ(userLoginRequest.User.Username)).Exists(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}
	if usernameExists {
		validationError := validate.Errors{
			Errors: map[string][]string{"username": {"exists"}},
		}
		response.FailedValidation(w, &validationError)
		return
	}

	emailExists, err := models.AppUsers(models.AppUserWhere.Email.EQ(userLoginRequest.User.Email)).Exists(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
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
		response.InternalServerError(w, err)
		return
	}

	newUser := models.AppUser{
		Username: userLoginRequest.User.Username,
		Password: hashedPassword,
		Email:    userLoginRequest.User.Email,
	}

	err = newUser.Insert(r.Context(), tx, boil.Infer())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	fmt.Println("userID", newUser.ID)
	token, done := app.createToken(w, r, newUser.ID)
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

func (app *application) usersGetCurrent(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	user, err := models.AppUsers(qm.Select(
		models.AppUserColumns.Email,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	var userResponse = dto.UserWrapper{
		User: dto.User{
			Email:    user.Email,
			Username: user.Username,
			Bio:      user.Bio.String,
			Image:    user.Image.String,
		},
	}

	response.JSON(w, http.StatusOK, userResponse)
}

func (app *application) usersUpdate(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	var userUpdateRequest dto.UserRequest
	err := request.DecodeJSON(w, r, &userUpdateRequest)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	_, err = models.AppUsers(qm.Select(models.AppUserColumns.ID), models.AppUserWhere.ID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	if userUpdateRequest.User.Username != "" {
		usernameExists, err := models.AppUsers(models.AppUserWhere.Username.EQ(userUpdateRequest.User.Username),
			models.AppUserWhere.ID.NEQ(userID)).Exists(r.Context(), tx)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		if usernameExists {
			validationError := validate.Errors{
				Errors: map[string][]string{"username": {"exists"}},
			}
			response.FailedValidation(w, &validationError)
			return
		}
	}

	if userUpdateRequest.User.Email != "" {
		emailExists, err := models.AppUsers(models.AppUserWhere.Email.EQ(userUpdateRequest.User.Email),
			models.AppUserWhere.ID.NEQ(userID)).Exists(r.Context(), tx)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		if emailExists {
			validationError := validate.Errors{
				Errors: map[string][]string{"email": {"exists"}},
			}
			response.FailedValidation(w, &validationError)
			return
		}
	}

	updates := models.M{}

	if userUpdateRequest.User.Username != "" {
		updates[models.AppUserColumns.Username] = userUpdateRequest.User.Username
	}
	if userUpdateRequest.User.Email != "" {
		updates[models.AppUserColumns.Email] = userUpdateRequest.User.Email
	}
	updates[models.AppUserColumns.Bio] = userUpdateRequest.User.Bio
	updates[models.AppUserColumns.Image] = userUpdateRequest.User.Image

	if userUpdateRequest.User.Password != "" {
		hashedPassword, err := argon2id.CreateHash(userUpdateRequest.User.Password, &argon2id.Params{
			Memory:      app.config.Argon2.Memory,
			Iterations:  app.config.Argon2.Iterations,
			Parallelism: app.config.Argon2.Parallelism,
			SaltLength:  app.config.Argon2.SaltLength,
			KeyLength:   app.config.Argon2.KeyLength,
		})
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		updates[models.AppUserColumns.Password] = hashedPassword
	}

	err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), tx, updates)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	updatedUser, err := models.AppUsers(qm.Select(
		models.AppUserColumns.Email,
		models.AppUserColumns.Username,
		models.AppUserColumns.Bio,
		models.AppUserColumns.Image),
		models.AppUserWhere.ID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	var userResponse = dto.UserWrapper{
		User: dto.User{
			Email:    updatedUser.Email,
			Username: updatedUser.Username,
			Bio:      updatedUser.Bio.String,
			Image:    updatedUser.Image.String,
		},
	}

	response.JSON(w, http.StatusOK, userResponse)
}
