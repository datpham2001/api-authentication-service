package controller

import (
	"mime/multipart"
	"realworld-authentication/dto/auth"
	"realworld-authentication/dto/user"
	"realworld-authentication/entity"
)

type AuthService interface {
	SignUp(input *auth.UserSignUpDto) (*entity.UserSignUpResponse, error)
	Login(input *auth.UserLoginDto) (*entity.UserLoginResponse, error)
	RefreshToken(input *auth.RefreshTokenRequestDto) (*entity.TokenResponse, error)
	Logout(userID string) error
	LoginWithGoogle(input *auth.GoogleLoginDto) (*entity.GoogleOauthTokenResponse, error)

	GetUserProfileByID(userID string) (*entity.UserProfileResponse, error)
	UpdateUserProfile(userID string, input *user.UserProfileUpdateDto) (*entity.UserProfileResponse, error)
	ResetPassword(userID string, input *user.UserResetPasswordDto) (*entity.UserPasswordResponse, error)
	ForgetPassword(email string) (*entity.UserPasswordResponse, error)
}

type FileService interface {
	UploadFile(fileName string, src multipart.File, fileType string) (*entity.UploadFileResponse, error)
	DeleteFile(fileName string) error
}
