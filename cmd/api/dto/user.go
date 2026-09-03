package dto

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type User struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserWrapper struct {
	User User `json:"user"`
}

type UserRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

type UserUpdateRequest struct {
	User struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"user"`
}

func ValidateUserLoginRequest(u *UserRequest) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "email",
			Field:   u.User.Email,
			Message: "required",
		},
		&validators.StringIsPresent{
			Name:    "password",
			Field:   u.User.Password,
			Message: "required",
		},
	)
}

func ValidateUserRegistrationRequest(u *UserRequest) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "email",
			Field:   u.User.Email,
			Message: "required",
		},
		&validators.StringIsPresent{
			Name:    "password",
			Field:   u.User.Password,
			Message: "required",
		},
		&validators.StringIsPresent{
			Name:    "username",
			Field:   u.User.Username,
			Message: "required",
		},
	)
}

func ValidateUserUpdateRequest(u *UserUpdateRequest) *validate.Errors {
	var fields []validate.Validator
	if u.User.Email != nil {
		fields = append(fields, &validators.StringIsPresent{
			Name: "email", Field: *u.User.Email, Message: "required",
		})
	}
	if u.User.Password != nil {
		fields = append(fields, &validators.StringIsPresent{
			Name: "password", Field: *u.User.Password, Message: "required",
		})
	}
	if u.User.Username != nil {
		fields = append(fields, &validators.StringIsPresent{
			Name: "username", Field: *u.User.Username, Message: "required",
		})
	}
	return validate.Validate(fields...)
}
