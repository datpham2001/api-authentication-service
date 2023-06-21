package entity

import "realworld-authentication/model"

type UserSignUpResponse struct {
	User *model.User `json:"user"`
}

func NewUserSignupResponse(u *model.User) *UserSignUpResponse {
	resp := new(UserSignUpResponse)
	resp.User = u
	return resp
}

type UserLoginResponse struct {
	User struct {
		Email       string `json:"email,omitempty"`
		Username    string `json:"username,omitempty"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

func NewUserLoginResponse(u *model.User) *UserLoginResponse {
	resp := new(UserLoginResponse)
	resp.User.Email = u.Email
	resp.User.Username = u.Username
	resp.User.AccessToken = u.AccessToken
	return resp
}

type TokenResponse struct {
	Token struct {
		AccessToken  string `json:"accessToken,omitempty"`
		RefreshToken string `json:"refreshToken,omitempty"`
	} `json:"token"`
}

func NewTokenResp(accessToken, refreshToken string) *TokenResponse {
	resp := new(TokenResponse)
	resp.Token.AccessToken = accessToken
	resp.Token.RefreshToken = refreshToken
	return resp
}

type GoogleOauthTokenResponse struct {
	Token struct {
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"token"`
}

func NewGoogleOauthTokenResp(accessToken string) *GoogleOauthTokenResponse {
	resp := new(GoogleOauthTokenResponse)
	resp.Token.AccessToken = accessToken
	return resp
}
