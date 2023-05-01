package controller

import (
	"fmt"
	"net/http"
	"realworld-authentication/config"
	"realworld-authentication/dto"
	"realworld-authentication/helper"
	"realworld-authentication/model/enum"
	"realworld-authentication/repository"
	"realworld-authentication/utils"

	"github.com/labstack/echo/v4"
)

type authController struct {
	Repository repository.AuthRepository
}

func NewAuthController(repo repository.AuthRepository) *authController {
	return &authController{
		Repository: repo,
	}
}

func (h *authController) SignUp(c echo.Context) error {
	var input dto.UserSignUpRequest

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Cannot parse input data. " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	switch input.User.Email {
	case "":
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "User email is not empty",
			ErrorCode: string(enum.ErrorCodeRequired.Email),
		})
	default:
		if !utils.ValidateEmail(input.User.Email) {
			return c.JSON(http.StatusBadRequest, &helper.APIResponse{
				Code:      http.StatusBadRequest,
				Status:    helper.APIStatus.Invalid,
				Message:   "User email is invalid format",
				ErrorCode: string(enum.ErrorCodeInvalid.Email),
			})
		}
	}

	switch input.User.Username {
	case "":
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Username is not empty",
			ErrorCode: string(enum.ErrorCodeRequired.Username),
		})
	default:
		if !utils.ValidateUsername(input.User.Username) {
			return c.JSON(http.StatusBadRequest, &helper.APIResponse{
				Code:      http.StatusBadRequest,
				Status:    helper.APIStatus.Invalid,
				Message:   "Username is invalid format",
				ErrorCode: string(enum.ErrorCodeInvalid.Username),
			})
		}
	}

	switch input.User.Password {
	case "":
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "User password is not empty",
			ErrorCode: string(enum.ErrorCodeRequired.Password),
		})
	default:
		if !utils.ValidatePassword(input.User.Password) {
			return c.JSON(http.StatusBadRequest, &helper.APIResponse{
				Code:      http.StatusBadRequest,
				Status:    helper.APIStatus.Invalid,
				Message:   "User password is invalid format",
				ErrorCode: string(enum.ErrorCodeInvalid.Password),
			})
		}
	}

	userSignupResponse := h.Repository.SignUp(&input)
	return c.JSON(userSignupResponse.Code, userSignupResponse)
}

func (h *authController) Login(c echo.Context) error {
	var input dto.UserLoginRequest

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Parse data error. " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	if input.User.Email == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Email must not be empty",
			ErrorCode: string(enum.ErrorCodeRequired.Email),
		})
	}

	if input.User.Password == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:      http.StatusBadRequest,
			Status:    helper.APIStatus.Invalid,
			Message:   "Password must not be empty",
			ErrorCode: string(enum.ErrorCodeRequired.Password),
		})
	}

	userLoginResp := h.Repository.Login(&input)
	return c.JSON(userLoginResp.Code, userLoginResp)
}

func (h *authController) RefreshToken(c echo.Context) error {
	var request dto.RefreshTokenRequest

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	refreshTokenResp := h.Repository.RefreshToken(&request)
	return c.JSON(refreshTokenResp.Code, refreshTokenResp)
}

func (h *authController) GoogleOauth(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
			Code:    http.StatusUnauthorized,
			Status:  helper.APIStatus.Unauthorized,
			Message: "Authorization code not be provided",
		})
	}

	pathUrl := "/"
	if c.QueryParam("state") != "" {
		pathUrl = c.QueryParam("state")
	}

	googleSignInResp := h.Repository.LoginWithGoogle(&dto.GoogleLoginRequest{
		AuthorizationCode: code,
		PathUrl:           pathUrl,
	})

	if googleSignInResp.Status != helper.APIStatus.Ok {
		return c.JSON(googleSignInResp.Code, googleSignInResp)
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    googleSignInResp.Data.(string),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   int(config.AppConfig.AccessTokenMaxAge) * 60,
	})

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(config.AppConfig.ClientOrigin, pathUrl))
}
