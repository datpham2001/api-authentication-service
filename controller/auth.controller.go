package controller

import (
	"fmt"
	"net/http"
	"realworld-authentication/config/env"
	"realworld-authentication/dto/auth"
	"realworld-authentication/dto/user"
	"realworld-authentication/helper"
	"realworld-authentication/model/enum"
	"realworld-authentication/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	AuthService AuthService
	FileService FileService
	Validator   *validator.Validate
}

func NewAuthController(authService AuthService, fileService FileService, validator *validator.Validate) *AuthController {
	return &AuthController{
		AuthService: authService,
		FileService: fileService,
		Validator:   validator,
	}
}

func (h *AuthController) SignUp(c echo.Context) error {
	var input auth.UserSignUpDto

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Bad request. Error: " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	err = h.Validator.Struct(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Validate error: " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.InvalidFields),
		})
	}

	if !utils.ValidateEmail(input.User.Email) {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "User email is invalid format",
			ErrorCode: string(enum.ErrorCodeInvalid.Email),
		})
	}

	if !utils.ValidateUsername(input.User.Username) {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Username is invalid format",
			ErrorCode: string(enum.ErrorCodeInvalid.Username),
		})
	}

	if !utils.ValidatePassword(input.User.Password) {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "User password is invalid format",
			ErrorCode: string(enum.ErrorCodeInvalid.Password),
		})
	}

	userSignupResponse, err := h.AuthService.SignUp(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Sign up successfully",
		Data:    userSignupResponse,
	})
}

func (h *AuthController) Login(c echo.Context) error {
	var input auth.UserLoginDto

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Parse data error. " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	err = h.Validator.Struct(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Validate error: " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.InvalidFields),
		})
	}

	userLoginResp, err := h.AuthService.Login(&input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Login successfully",
		Data:    userLoginResp,
	})
}

func (h *AuthController) RefreshToken(c echo.Context) error {
	var request auth.RefreshTokenRequestDto

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	refreshTokenResp, err := h.AuthService.RefreshToken(&request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Refresh token succcessfully",
		Data:    refreshTokenResp,
	})
}

func (h *AuthController) GoogleOauth(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
			Status:  helper.APIStatus.Unauthorized,
			Message: "Authorization code not be provided",
		})
	}

	pathUrl := "/"
	if c.QueryParam("state") != "" {
		pathUrl = c.QueryParam("state")
	}

	googleSignInResp, err := h.AuthService.LoginWithGoogle(&auth.GoogleLoginDto{
		AuthorizationCode: code,
		PathUrl:           pathUrl,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    googleSignInResp.Token.AccessToken,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   int(env.AppConfig.AccessTokenMaxAge) * 60,
	})

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(env.AppConfig.ClientOrigin, pathUrl))
}

func (h *AuthController) GetUserProfileByID(c echo.Context) error {
	var (
		userID = c.Param("userID")
	)

	userProfileResp, err := h.AuthService.GetUserProfileByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Get user profile successfully",
		Data:    userProfileResp,
	})
}

func (h *AuthController) UpdateUserProfile(c echo.Context) error {
	var (
		userID = c.Param("userID")
		input  user.UserProfileUpdateDto
	)

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Parse data error. " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	err = h.Validator.Struct(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Validate error: " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.InvalidFields),
		})
	}

	userUpdateProfileResp, err := h.AuthService.UpdateUserProfile(userID, &input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Update user profile successfully",
		Data:    userUpdateProfileResp,
	})
}

func (h *AuthController) GetMyProfile(c echo.Context) error {
	var (
		userID = getUserIDFromToken(c)
	)

	if userID == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: "Missing User ID",
		})
	}

	myProfileResp, err := h.AuthService.GetUserProfileByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Get user profile successfully",
		Data:    myProfileResp,
	})
}

func (h *AuthController) ResetUserPassword(c echo.Context) error {
	var (
		userID = getUserIDFromToken(c)
		input  user.UserResetPasswordDto
	)

	if userID == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: "Missing User ID",
		})
	}

	err := c.Bind(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Parse data error. " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.ParseData),
		})
	}

	err = h.Validator.Struct(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:    helper.APIStatus.Invalid,
			Message:   "Validate error: " + err.Error(),
			ErrorCode: string(enum.ErrorCodeInvalid.InvalidFields),
		})
	}

	if !utils.ValidatePassword(input.User.NewPassword) {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: "Password is invalid format",
		})
	}

	userResetPassword, err := h.AuthService.ResetPassword(userID, &input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Reset user password successfully",
		Data:    userResetPassword,
	})
}

func (h *AuthController) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: err.Error(),
		})
	}
	fileType := file.Header.Get("Content-Type")

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadGateway, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}
	defer src.Close()

	uploadFileResp, err := h.FileService.UploadFile(file.Filename, src, fileType)
	if err != nil {
		return c.JSON(http.StatusBadGateway, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Upload file successfully",
		Data:    uploadFileResp,
	})
}

func (h *AuthController) Logout(c echo.Context) error {
	var (
		userID = getUserIDFromToken(c)
	)

	if userID == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: "Missing User ID",
		})
	}

	err := h.AuthService.Logout(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Logout successfully",
	})
}

func (h *AuthController) ForgetPassword(c echo.Context) error {
	var userEmail = c.QueryParam("email")
	if userEmail == "" {
		return c.JSON(http.StatusBadRequest, &helper.APIResponse{
			Status:  helper.APIStatus.Invalid,
			Message: "Missing email",
		})
	}

	forgetPasswordResp, err := h.AuthService.ForgetPassword(userEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &helper.APIResponse{
			Status:  helper.APIStatus.Error,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &helper.APIResponse{
		Status:  helper.APIStatus.Ok,
		Message: "Reset password successfully",
		Data:    forgetPasswordResp,
	})
}

func getUserIDFromToken(c echo.Context) string {
	userID, ok := c.Get("userId").(string)
	if !ok {
		return ""
	}

	return userID
}
