package middleware

import (
	"errors"
	"net/http"
	"realworld-authentication/config"
	"realworld-authentication/helper"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header value
		token, err := extractTokenFromHeaderString(c.Request().Header.Get("Authorization"))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		claims, err := helper.ValidateToken(token, config.AppConfig.AccessTokenPublicKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
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
