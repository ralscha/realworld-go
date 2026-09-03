package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/alexedwards/argon2id"
	"github.com/gobuffalo/validate"
	"realworldgo.rasc.ch/cmd/api/dto"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/internal/models"
	"realworldgo.rasc.ch/internal/request"
	"realworldgo.rasc.ch/internal/response"
)

var userNotFoundPasswordHash string

const userNotFoundPassword string = "userNotFoundPassword"

func initAuth(config config.Config) error {
	var err error
	userNotFoundPasswordHash, err = argon2id.CreateHash(userNotFoundPassword, passwordHashParams(config))
	return err
}

func passwordHashParams(config config.Config) *argon2id.Params {
	return &argon2id.Params{
		Memory:      config.Argon2.Memory,
		Iterations:  config.Argon2.Iterations,
		Parallelism: config.Argon2.Parallelism,
		SaltLength:  config.Argon2.SaltLength,
		KeyLength:   config.Argon2.KeyLength,
	}
}

func userResponse(user *models.AppUser, token string) dto.UserWrapper {
	return dto.UserWrapper{User: dto.User{
		Email:    user.Email,
		Token:    token,
		Username: user.Username,
		Bio:      user.Bio.String,
		Image:    user.Image.String,
	}}
}

func (app *application) usersLogin(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)

	var userLoginRequest dto.UserRequest
	if ok := request.DecodeJSONValidate(w, r, &userLoginRequest, dto.ValidateUserLoginRequest); !ok {
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
			if _, err := argon2id.ComparePasswordAndHash(userNotFoundPassword, userNotFoundPasswordHash); err != nil {
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

	token, err := app.createToken(r, user.ID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, userResponse(user, token))
}

func (app *application) createToken(r *http.Request, userID int64) (string, error) {
	ctx, err := app.sessionManager.Load(r.Context(), "")
	if err != nil {
		return "", err
	}
	app.sessionManager.Put(ctx, "userID", userID)
	token, _, err := app.sessionManager.Commit(ctx)
	if err != nil {
		return "", err
	}
	return token, nil
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

	hashedPassword, err := argon2id.CreateHash(userLoginRequest.User.Password, passwordHashParams(*app.config))
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

	token, err := app.createToken(r, newUser.ID)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, userResponse(&newUser, token))
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

	response.JSON(w, http.StatusOK, userResponse(user, app.sessionManager.Token(r.Context())))
}

func (app *application) usersUpdate(w http.ResponseWriter, r *http.Request) {
	tx := r.Context().Value(transactionKey).(*sql.Tx)
	var userUpdateRequest dto.UserUpdateRequest
	if ok := request.DecodeJSONValidate(w, r, &userUpdateRequest, dto.ValidateUserUpdateRequest); !ok {
		return
	}

	userID := app.sessionManager.GetInt64(r.Context(), "userID")
	_, err := models.AppUsers(qm.Select(models.AppUserColumns.ID), models.AppUserWhere.ID.EQ(userID)).One(r.Context(), tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	if userUpdateRequest.User.Username != nil {
		usernameExists, err := models.AppUsers(models.AppUserWhere.Username.EQ(*userUpdateRequest.User.Username),
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

	if userUpdateRequest.User.Email != nil {
		emailExists, err := models.AppUsers(models.AppUserWhere.Email.EQ(*userUpdateRequest.User.Email),
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

	if userUpdateRequest.User.Username != nil {
		updates[models.AppUserColumns.Username] = *userUpdateRequest.User.Username
	}
	if userUpdateRequest.User.Email != nil {
		updates[models.AppUserColumns.Email] = *userUpdateRequest.User.Email
	}
	if userUpdateRequest.User.Bio != nil {
		updates[models.AppUserColumns.Bio] = *userUpdateRequest.User.Bio
	}
	if userUpdateRequest.User.Image != nil {
		updates[models.AppUserColumns.Image] = *userUpdateRequest.User.Image
	}

	if userUpdateRequest.User.Password != nil {
		hashedPassword, err := argon2id.CreateHash(*userUpdateRequest.User.Password, passwordHashParams(*app.config))
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
		updates[models.AppUserColumns.Password] = hashedPassword
	}

	if len(updates) > 0 {
		err = models.AppUsers(models.AppUserWhere.ID.EQ(userID)).UpdateAll(r.Context(), tx, updates)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}
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

	response.JSON(w, http.StatusOK, userResponse(updatedUser, app.sessionManager.Token(r.Context())))
}
