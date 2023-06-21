package middleware

import (
	"errors"
	"net/http"
	"realworld-authentication/config/env"
	"realworld-authentication/helper"
	"realworld-authentication/model/enum"
	auth_service "realworld-authentication/service/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	authStorage auth_service.AuthStorage
}

func NewAuthMiddleware(s auth_service.AuthStorage) *AuthMiddleware {
	return &AuthMiddleware{
		authStorage: s,
	}
}

func (m *AuthMiddleware) TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header value
		token, err := extractTokenFromHeaderString(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
				Status:  helper.APIStatus.Unauthorized,
				Message: err.Error(),
			})
		}

		claims, err := helper.ValidateToken(token, env.AppConfig.AccessTokenKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
				Status:  helper.APIStatus.Unauthorized,
				Message: err.Error(),
			})
		}

		c.Set("userId", claims.UserID)
		return next(c)
	}
}

func (m *AuthMiddleware) RolePermissionAuthorize(role enum.UserRoleValue, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header value
		token, err := extractTokenFromHeaderString(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
				Status:  helper.APIStatus.Unauthorized,
				Message: err.Error(),
			})
		}

		claims, err := helper.ValidateToken(token, env.AppConfig.AccessTokenKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &helper.APIResponse{
				Status:  helper.APIStatus.Unauthorized,
				Message: err.Error(),
			})
		}

		// get role from userId to check
		user, err := m.authStorage.GetUserByID(claims.UserID)
		if err != nil || user.Role != role {
			return c.JSON(http.StatusForbidden, &helper.APIResponse{
				Status:  helper.APIStatus.Invalid,
				Message: "Your account cannot perform this action",
			})
		}

		c.Set("userId", claims.UserID)
		return next(c)
	}
}

func extractTokenFromHeaderString(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) < 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("missing access token")
	}

	return parts[1], nil
}
