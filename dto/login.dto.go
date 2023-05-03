package dto

type UserLoginRequest struct {
	User struct {
		Email    string `json:"email,omitempty" validate:"required"`
		Password string `json:"password,omitempty" validate:"required"`
	} `json:"user" validate:"required"`
}

type GoogleLoginRequest struct {
	AuthorizationCode string
	PathUrl           string
}
