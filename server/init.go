package server

import (
	"fmt"
	"os"
	"realworld-authentication/controller"
	auth_middleware "realworld-authentication/middleware"
	"realworld-authentication/repository"
	auth_service "realworld-authentication/service/auth"
	file_service "realworld-authentication/service/file"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPServer struct {
	Router         *echo.Echo
	Validator      *validator.Validate
	AuthMiddlware  *auth_middleware.AuthMiddleware
	FileStorage    file_service.FileStorage
	FileService    controller.FileService
	AuthStorage    auth_service.AuthStorage
	AuthService    controller.AuthService
	AuthController *controller.AuthController
}

func (server *HTTPServer) Init(db *mongo.Database) {
	server.Router = echo.New()
	server.Validator = validator.New()
	server.AuthStorage = repository.NewAuthStorage(db)
	server.FileStorage = repository.NewFileStorage(db)
	server.AuthMiddlware = auth_middleware.NewAuthMiddleware(server.AuthStorage)
	server.FileService = file_service.NewFileService(server.FileStorage)
	server.AuthService = auth_service.NewAuthService(server.AuthStorage, server.FileService)
	server.AuthController = controller.NewAuthController(server.AuthService, server.FileService, server.Validator)
}

func (server *HTTPServer) UseMiddleware() {
	server.Router.Pre(middleware.RemoveTrailingSlash())

	logger := zerolog.New(os.Stdout)
	server.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("uri", v.URI).
				Int("status", v.Status).
				Int64("responseSize", v.ResponseSize).
				Msg("request")
			return nil
		},
	}))

	server.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.HEAD},
	}))
}

func (server *HTTPServer) Launch(port int64) {
	fmt.Printf("Listening on port %d \n", port)
	server.Router.Logger.Fatal(server.Router.Start(fmt.Sprintf(":%d", port)))
}
