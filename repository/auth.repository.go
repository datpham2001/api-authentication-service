package repository

import (
	"net/http"
	"realworld-authentication/config"
	"realworld-authentication/dto"
	"realworld-authentication/entity"
	"realworld-authentication/helper"
	"realworld-authentication/model"
	"realworld-authentication/model/enum"
	"realworld-authentication/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthRepository(db *mongo.Database) AuthRepository {
	r := &Instance{
		ColName:        "authentication",
		TemplateObject: &model.User{},
	}

	r.ApplyDatabase(db)
	return r
}

func (r *Instance) SignUp(input *dto.UserSignUpRequest) *helper.APIResponse {
	user := &model.User{
		Email:    input.User.Email,
		Username: input.User.Username,
	}

	resp := r.QueryOne(model.User{
		ComplexQuery: []*bson.M{
			{
				"$or": []*bson.M{{
					"username": user.Username,
				}, {
					"email": user.Email,
				}},
			},
		},
	})

	if resp.Status == helper.APIStatus.Ok {
		return &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Username or email is existed",
			ErrorCode: string(enum.ErrorCodeExisted.UsernameOrEmail),
		}
	}

	hashedPassword, err := helper.HashPassword(input.User.Password)
	if err != nil {
		return &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Cannot hash user password " + err.Error(),
			ErrorCode: string(enum.ErrorCodePackage.Bcrypt),
		}
	}

	user.HashedPassword = hashedPassword
	user.UserID = utils.GenAccountID()
	user.Status = enum.UserStatus.Active
	user.Role = enum.UserRole.User
	user.Role = enum.UserRole.User

	userCreateResp := r.Create(user)
	if userCreateResp.Status != helper.APIStatus.Ok {
		return userCreateResp
	}

	return &helper.APIResponse{
		Code:    http.StatusCreated,
		Status:  helper.APIStatus.Ok,
		Message: "Signup user successfully",
		Data:    userCreateResp.Data.([]*model.User)[0],
	}
}

func (r *Instance) Login(input *dto.UserLoginRequest) *helper.APIResponse {
	existUserResp := r.QueryOne(model.User{
		Email: input.User.Email,
	})
	if existUserResp.Status != helper.APIStatus.Ok {
		return existUserResp
	}

	user := existUserResp.Data.([]*model.User)[0]
	if !helper.VerifyPassword(user.HashedPassword, input.User.Password) {
		return &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Password is not matched",
			ErrorCode: string(enum.ErrorCodeInvalid.Password),
		}
	}

	accessToken, err := helper.GenerateJWT(user.UserID, config.AppConfig.AccessTokenExpiredIn, config.AppConfig.AccessTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  helper.APIStatus.Error,
			Message: "Failed to generate token",
		}
	}
	user.AccessToken = *accessToken.Token

	refreshToken, err := helper.GenerateJWT(user.UserID, config.AppConfig.RefreshTokenExpiredIn, config.AppConfig.RefreshTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  helper.APIStatus.Error,
			Message: "Failed to generate token",
		}
	}
	user.RefreshToken = *refreshToken.Token

	resp := r.UpdateOne(model.User{ID: user.ID}, user)
	if resp.Status != helper.APIStatus.Ok {
		return resp
	}

	return &helper.APIResponse{
		Code:    http.StatusOK,
		Status:  helper.APIStatus.Ok,
		Message: "Login user successfully",
		Data:    entity.NewUserLoginResp(user),
	}
}

func (r *Instance) RefreshToken(input *dto.RefreshTokenRequest) *helper.APIResponse {
	token, err := helper.ValidateToken(input.RefreshToken, config.AppConfig.RefreshTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		}
	}

	resp := r.QueryOne(model.User{UserID: token.UserID})
	if resp.Status != helper.APIStatus.Ok {
		return resp
	}

	accessToken, err := helper.GenerateJWT(token.UserID, config.AppConfig.AccessTokenExpiredIn, config.AppConfig.AccessTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		}
	}

	refreshToken, err := helper.GenerateJWT(token.UserID, config.AppConfig.RefreshTokenExpiredIn, config.AppConfig.RefreshTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		}
	}

	// update new refresh token in db for the user
	resp = r.UpdateOne(model.User{
		UserID: token.UserID,
	}, model.User{
		RefreshToken: *refreshToken.Token,
	})
	if resp.Status != helper.APIStatus.Ok {
		return resp
	}

	return &helper.APIResponse{
		Code:    http.StatusOK,
		Status:  helper.APIStatus.Ok,
		Message: "Refresh token successfully",
		Data:    entity.NewTokenResp(*accessToken.Token, *refreshToken.Token),
	}
}

func (r *Instance) LoginWithGoogle(input *dto.GoogleLoginRequest) *helper.APIResponse {
	tokenResp, err := helper.GetGoogleOauthToken(input.AuthorizationCode)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusBadGateway,
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		}
	}

	googleUserInfo, err := helper.GetGoogleUserInfo(tokenResp.AccessToken, tokenResp.TokenID)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusBadGateway,
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		}
	}

	var (
		userEmail = strings.ToLower(googleUserInfo.Email)
		userData  *model.User
	)

	userResp := r.QueryOne(model.User{Email: userEmail})
	if userResp.Status == helper.APIStatus.Notfound {
		userData = &model.User{
			UserID:   googleUserInfo.ID,
			Email:    userEmail,
			Username: googleUserInfo.Name,
			Provider: enum.ProviderName.Google,
		}

		resp := r.Create(userData)
		if resp.Status != helper.APIStatus.Ok {
			return resp
		}
	} else if userResp.Status == helper.APIStatus.Ok {
		userData = userResp.Data.([]*model.User)[0]
	}

	accessToken, err := helper.GenerateJWT(userData.UserID, config.AppConfig.AccessTokenExpiredIn, config.AppConfig.AccessTokenKey)
	if err != nil {
		return &helper.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		}
	}

	userData.AccessToken = *accessToken.Token
	return &helper.APIResponse{
		Code:   http.StatusOK,
		Status: helper.APIStatus.Ok,
		Data:   *accessToken.Token,
	}
}
