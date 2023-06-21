package auth

type UserLoginDto struct {
	User struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	} `json:"user" validate:"required"`
}

type GoogleLoginDto struct {
	AuthorizationCode string
	PathUrl           string
}
