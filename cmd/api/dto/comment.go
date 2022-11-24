package dto

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Comment struct {
	Id        int64   `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Body      string  `json:"body"`
	Author    Profile `json:"author"`
}

type CommentOne struct {
	Comment Comment `json:"comment"`
}

type CommentMany struct {
	Comments []Comment `json:"comments"`
}

type CommentRequest struct {
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

func ValidateCommentRequest(c *CommentRequest) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "body",
			Field:   c.Comment.Body,
			Message: "required",
		},
	)
}
