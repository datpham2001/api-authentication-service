package main

import (
	"net/http"
	"realworld-authentication/config/db"
	"realworld-authentication/config/env"
	"realworld-authentication/helper"
	"realworld-authentication/server"

	"github.com/labstack/echo/v4"
)

var app *server.HTTPServer = &server.HTTPServer{}

func init() {
	// load configuration
	err := env.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// connect database
	db.ConnectDB()

	// setup server app
	app.Init(db.Client.Database(env.AppConfig.DBName))
}

func main() {
	app.UseMiddleware()

	// setup API server
	app.Router.GET("/api/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &helper.APIResponse{
			Status:  helper.APIStatus.Ok,
			Message: "Health check api successfully",
		})
	})

	// auth route
	{
		app.Router.POST("/api/auth/signup", app.AuthController.SignUp)
		app.Router.POST("/api/auth/login", app.AuthController.Login)
		app.Router.POST("/api/auth/token/refresh", app.AuthController.RefreshToken)
		app.Router.POST("/api/auth/logout", app.AuthController.Logout, app.AuthMiddlware.TokenAuthMiddleware)
		app.Router.GET("/api/sessions/oauth/google", app.AuthController.GoogleOauth, app.AuthMiddlware.TokenAuthMiddleware)
	}

	// user route
	{
		app.Router.GET("/api/users/:userID/profile", app.AuthController.GetUserProfileByID)
		app.Router.GET("/api/users/me/profile", app.AuthController.GetMyProfile, app.AuthMiddlware.TokenAuthMiddleware)
		app.Router.PUT("/api/users/:userID/profile", app.AuthController.UpdateUserProfile, app.AuthMiddlware.TokenAuthMiddleware)
		app.Router.PUT("/api/users/me/reset-password", app.AuthController.ResetUserPassword, app.AuthMiddlware.TokenAuthMiddleware)
		app.Router.PUT("/api/users/forget-password", app.AuthController.ForgetPassword)
		app.Router.POST("/api/upload", app.AuthController.UploadFile, app.AuthMiddlware.TokenAuthMiddleware)
	}

	// launch app
	app.Launch(env.AppConfig.Port)
}
