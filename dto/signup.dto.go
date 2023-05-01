package dto

type UserSignUpRequest struct {
	User struct {
		Email    string `json:"email,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	} `json:"user"`
}
