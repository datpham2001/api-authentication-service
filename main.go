package main

import (
	"fmt"
	"net/http"
	"realworld-authentication/config"
	"realworld-authentication/controller"
	"realworld-authentication/middleware"
	"realworld-authentication/repository"
	"realworld-authentication/route"
	"realworld-authentication/storage/db"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var server *echo.Echo

func init() {
	// load configuration
	err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// connect database
	db.ConnectDB()

	// setup server
	server = route.New()
}

func main() {
	// init instances
	validator := validator.New()
	repository := repository.NewAuthRepository(db.Client.Database(config.AppConfig.DBName))
	controller := controller.NewAuthController(repository, validator)

	// setup api server & its middleware
	router := server.Group("/api")
	router.GET("/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Start authentication server")
	})

	authRouter := router.Group("/auth")

	authRouter.POST("/signup", controller.SignUp)
	authRouter.POST("/login", controller.Login)
	authRouter.POST("/token/refresh", controller.RefreshToken, middleware.AuthMiddleware)

	router.GET("/sessions/oauth/google", controller.GoogleOauth, middleware.AuthMiddleware)

	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", config.AppConfig.Port)))
}
