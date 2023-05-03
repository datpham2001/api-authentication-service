package dto

type UserSignUpRequest struct {
	User struct {
		Email    string `json:"email" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user" validate:"required"`
}
