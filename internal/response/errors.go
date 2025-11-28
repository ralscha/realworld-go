// Package response provides HTTP response utilities.
package response

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gobuffalo/validate"
)

func errorMessage(w http.ResponseWriter, status int, message string) {
	JSONWithHeaders(w, status, map[string]string{"error": message}, nil)
}

func InternalServerError(w http.ResponseWriter, err error) {
	slog.Error(err.Error(), "error", err)

	message := "The server encountered a problem and could not process your request"
	errorMessage(w, http.StatusInternalServerError, message)
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	message := "The requested resource could not be found"
	errorMessage(w, http.StatusNotFound, message)
}

func Unauthorized(w http.ResponseWriter) {
	message := "Authentication is required to access this resource"
	errorMessage(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter) {
	message := "You are not authorized to access this resource"
	errorMessage(w, http.StatusForbidden, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	errorMessage(w, http.StatusMethodNotAllowed, message)
}

func BadRequest(w http.ResponseWriter, err error) {
	errorMessage(w, http.StatusBadRequest, err.Error())
}

func FailedValidation(w http.ResponseWriter, v *validate.Errors) {
	JSON(w, http.StatusUnprocessableEntity, v)
}
