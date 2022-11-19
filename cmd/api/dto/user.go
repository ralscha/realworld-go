package dto

type UserWrapper struct {
	User struct {
		Email    string      `json:"email"`
		Token    string      `json:"token"`
		Username string      `json:"username"`
		Bio      string      `json:"bio"`
		Image    interface{} `json:"image"`
	} `json:"user"`
}
