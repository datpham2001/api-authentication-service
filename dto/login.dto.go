package dto

type UserLoginRequest struct {
	User struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"user"`
}

type GoogleLoginRequest struct {
	AuthorizationCode string
	PathUrl           string
}
