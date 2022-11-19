package dto

import "time"

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type Comment struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Author    `json:"author"`
}

type CommentOne struct {
	Comment Comment `json:"comment"`
}

type CommentMany struct {
	Comments []Comment `json:"comments"`
}
