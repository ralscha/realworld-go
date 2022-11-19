package dto

import "github.com/gobuffalo/validate"

type Validatable interface {
	Validate() *validate.Errors
}
