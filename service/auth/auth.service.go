package auth

import (
	"errors"
	"realworld-authentication/config/env"
	"realworld-authentication/controller"
	"realworld-authentication/dto/auth"
	"realworld-authentication/dto/user"
	"realworld-authentication/entity"
	"realworld-authentication/helper"
	"realworld-authentication/model"
	"realworld-authentication/model/enum"
	"realworld-authentication/utils"
	"strings"
)

type authService struct {
	storage     AuthStorage
	fileService controller.FileService
}

func NewAuthService(storage AuthStorage, fileService controller.FileService) *authService {
	return &authService{
		storage:     storage,
		fileService: fileService,
	}
}

func (s *authService) SignUp(input *auth.UserSignUpDto) (*entity.UserSignUpResponse, error) {
	user := &model.User{
		Email:    input.User.Email,
		Username: input.User.Username,
	}

	_, err := s.storage.GetUserByUsernameOrEmail(user.Username, user.Email)
	if err == nil {
		return nil, errors.New("username or email is existed")
	}

	hashedPassword, err := helper.HashPassword(input.User.Password)
	if err != nil {
		return nil, err
	}
	user.HashedPassword = hashedPassword

	user.UserID = utils.GenAccountID()
	user.Status = enum.UserStatus.Active
	user.Role = enum.UserRole.User

	userCreateResp, err := s.storage.CreateUser(user)
	if err != nil {
		return nil, err
	}

	userSignupEntity := entity.NewUserSignupResponse(userCreateResp)
	return userSignupEntity, nil
}

func (s *authService) Login(input *auth.UserLoginDto) (*entity.UserLoginResponse, error) {
	existUserResp, err := s.storage.GetUserByEmail(input.User.Email)
	if err != nil {
		return nil, err
	}

	if !helper.VerifyPassword(existUserResp.HashedPassword, input.User.Password) {
		return nil, errors.New("password is not matched")
	}

	accessToken, err := helper.GenerateJWT(existUserResp.UserID, env.AppConfig.AccessTokenExpiredIn, env.AppConfig.AccessTokenKey)
	if err != nil {
		return nil, err
	}
	existUserResp.AccessToken = *accessToken.Token

	refreshToken, err := helper.GenerateJWT(existUserResp.UserID, env.AppConfig.RefreshTokenExpiredIn, env.AppConfig.RefreshTokenKey)
	if err != nil {
		return nil, err
	}
	existUserResp.RefreshToken = *refreshToken.Token

	_, err = s.storage.UpdateUser(&model.User{ID: existUserResp.ID}, existUserResp)
	if err != nil {
		return nil, err
	}

	return entity.NewUserLoginResponse(existUserResp), nil
}

func (s *authService) RefreshToken(input *auth.RefreshTokenRequestDto) (*entity.TokenResponse, error) {
	token, err := helper.ValidateToken(input.RefreshToken, env.AppConfig.RefreshTokenKey)
	if err != nil {
		return nil, err
	}

	_, err = s.storage.GetUserByID(token.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := helper.GenerateJWT(token.UserID, env.AppConfig.AccessTokenExpiredIn, env.AppConfig.AccessTokenKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateJWT(token.UserID, env.AppConfig.RefreshTokenExpiredIn, env.AppConfig.RefreshTokenKey)
	if err != nil {
		return nil, err
	}

	// update new refresh token in db for the user
	_, err = s.storage.UpdateUser(&model.User{
		UserID: token.UserID,
	}, &model.User{
		RefreshToken: *refreshToken.Token,
	})
	if err != nil {
		return nil, err
	}

	return entity.NewTokenResp(*accessToken.Token, *refreshToken.Token), nil
}

func (s *authService) LoginWithGoogle(input *auth.GoogleLoginDto) (*entity.GoogleOauthTokenResponse, error) {
	tokenResp, err := helper.GetGoogleOauthToken(input.AuthorizationCode)
	if err != nil {
		return nil, err
	}

	googleUserInfo, err := helper.GetGoogleUserInfo(tokenResp.AccessToken, tokenResp.TokenID)
	if err != nil {
		return nil, err
	}

	var (
		userEmail = strings.ToLower(googleUserInfo.Email)
	)

	userResp, err := s.storage.GetUserByEmail(userEmail)
	if err != nil {
		userResp = &model.User{
			UserID:   googleUserInfo.ID,
			Email:    userEmail,
			Username: googleUserInfo.Name,
			Provider: enum.ProviderName.Google,
			Role:     enum.UserRole.User,
		}

		_, err := s.storage.CreateUser(userResp)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := helper.GenerateJWT(userResp.UserID, env.AppConfig.AccessTokenExpiredIn, env.AppConfig.AccessTokenKey)
	if err != nil {
		return nil, err
	}

	userResp.AccessToken = *accessToken.Token
	return entity.NewGoogleOauthTokenResp(userResp.AccessToken), nil
}

func (s *authService) GetUserProfileByID(userID string) (*entity.UserProfileResponse, error) {
	resp, err := s.storage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return entity.NewUserProfileResponse(resp), nil
}

func (s *authService) UpdateUserProfile(userID string, input *user.UserProfileUpdateDto) (*entity.UserProfileResponse, error) {
	existUser, err := s.storage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	updateData := &model.User{}
	if input.User.Email != "" {
		updateData.Email = input.User.Email
	}
	if input.User.Username != "" {
		updateData.Username = input.User.Username
	}
	if input.User.Bio != nil && *input.User.Bio != "" {
		updateData.Bio = input.User.Bio
	}
	if input.User.Avatar != nil {
		updateData.Avatar = input.User.Avatar
	}

	updateUserResp, err := s.storage.UpdateUser(&model.User{
		ID: existUser.ID,
	}, updateData)
	if err != nil {
		return nil, err
	}

	return entity.NewUserProfileResponse(updateUserResp), nil
}

func (s *authService) ResetPassword(userID string, input *user.UserResetPasswordDto) (*entity.UserPasswordResponse, error) {
	existUser, err := s.storage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if !helper.VerifyPassword(existUser.HashedPassword, input.User.CurrentPassword) {
		return nil, errors.New("current password is not matched")
	}

	hashedPassword, err := helper.HashPassword(input.User.NewPassword)
	if err != nil {
		return nil, err
	}

	updateUserPassword, err := s.storage.UpdateUserPassword(&model.User{
		ID: existUser.ID,
	}, hashedPassword)
	if err != nil {
		return nil, err
	}

	return entity.NewUserPasswordResponse(updateUserPassword), nil
}

func (s *authService) Logout(userID string) error {
	existUser, err := s.storage.GetUserByID(userID)
	if err != nil {
		return err
	}

	// revoke user refresh token in db
	return s.storage.DeleteToken(existUser.RefreshToken)
}

func (s *authService) ForgetPassword(email string) (*entity.UserPasswordResponse, error) {
	existUser, err := s.storage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// send email to get new password
	randomToken := utils.RandomString(utils.PASSWORD_LENGTH)
	err := client.NotificationService.SendEmail(email, randomToken)
	if err != nil {
		return nil, err
	}

	updateUserResp, err = s.storage.UpdateUserPassword(&model.User{
		ID: existUser.ID,
	}, randomToken)
	if err != nil {
		return nil, err
	}

	return entity.NewUserPasswordResponse(updateUserResp), nil
}
