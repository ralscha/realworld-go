package dto

import (
	"github.com/gobuffalo/validate"
)

type ValidatorFn[T any] func(o T) *validate.Errors
