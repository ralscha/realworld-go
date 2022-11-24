package dto

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Article struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         Profile  `json:"author"`
}

type ArticleOne struct {
	Article Article `json:"article"`
}

type ArticlesMany struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int64     `json:"articlesCount"`
}

type ArticleRequest struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

func ValidateArticleCreateRequest(a *ArticleRequest) *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{
			Name:    "title",
			Field:   a.Article.Title,
			Message: "required",
		},
		&validators.StringIsPresent{
			Name:    "description",
			Field:   a.Article.Description,
			Message: "required",
		},
		&validators.StringIsPresent{
			Name:    "body",
			Field:   a.Article.Body,
			Message: "required",
		},
	)
}
