package repository

import (
	"realworld-authentication/dto"
	"realworld-authentication/helper"
)

type AuthRepository interface {
	SignUp(input *dto.UserSignUpRequest) *helper.APIResponse
	Login(input *dto.UserLoginRequest) *helper.APIResponse
	RefreshToken(input *dto.RefreshTokenRequest) *helper.APIResponse
	LoginWithGoogle(input *dto.GoogleLoginRequest) *helper.APIResponse
}
