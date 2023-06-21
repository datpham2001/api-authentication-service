package entity

import (
	"realworld-authentication/model"
	"realworld-authentication/model/enum"
)

type UserProfileResponse struct {
	User *model.User `json:"user"`
}

func NewUserProfileResponse(u *model.User) *UserProfileResponse {
	resp := new(UserProfileResponse)
	resp.User = u

	return resp
}

type UserPasswordResponse struct {
	User struct {
		Email             string `json:"email,omitempty"`
		Username          string `json:"username,omitempty"`
		IsChangedPassword *bool  `json:"isChangedPassword,omitempty"`
	} `json:"user"`
}

func NewUserPasswordResponse(u *model.User) *UserPasswordResponse {
	resp := new(UserPasswordResponse)
	resp.User.Email = u.Email
	resp.User.Username = u.Username
	resp.User.IsChangedPassword = &enum.TRUE

	return resp
}
