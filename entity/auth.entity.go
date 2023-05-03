package entity

import "realworld-authentication/model"

type UserSignUpResponse struct {
	User *model.User `json:"user"`
}

type userLoginResponse struct {
	User struct {
		Email       string `json:"email,omitempty"`
		Username    string `json:"username,omitempty"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

func NewUserLoginResp(u *model.User) *userLoginResponse {
	resp := new(userLoginResponse)
	resp.User.Email = u.Email
	resp.User.Username = u.Username
	resp.User.AccessToken = u.AccessToken
	return resp
}

type tokenResponse struct {
	Token struct {
		AccessToken  string `json:"accessToken,omitempty"`
		RefreshToken string `json:"refreshToken,omitempty"`
	} `json:"token"`
}

func NewTokenResp(accessToken, refreshToken string) *tokenResponse {
	resp := new(tokenResponse)
	resp.Token.AccessToken = accessToken
	resp.Token.RefreshToken = refreshToken
	return resp
}
