package dto

import "time"

type Comment struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Profile   `json:"author"`
}

type CommentOne struct {
	Comment Comment `json:"comment"`
}

type CommentMany struct {
	Comments []Comment `json:"comments"`
}
