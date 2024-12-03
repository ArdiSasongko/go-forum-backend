package v1

import (
	userhandler "github.com/ArdiSasongko/go-forum-backend/internal/handler/user.handler"
	userservice "github.com/ArdiSasongko/go-forum-backend/internal/service/user.service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type ApiRouter struct {
	userService userservice.UserService
}

func NewApiRouter(userService userservice.UserService) *ApiRouter {
	return &ApiRouter{userService: userService}
}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	authGroup := app.Group("/user")
	authGroupV1 := authGroup.Group("/v1")

	userHandler := userhandler.NewUserHandler(h.userService)
	authGroupV1.Post("/register", userHandler.Register)
	authGroupV1.Post("/login", userHandler.Login)
	authGroupV1.Put("/refresh-token", middleware.MiddlewareRefreshToken, userHandler.RefreshToken)
	authGroupV1.Put("/validate/email", middleware.MiddlewareAuthValidate, userHandler.ValidateUser)
	authGroupV1.Get("/validate/resend", middleware.MiddlewareAuthValidate, userHandler.ResendEmail)
	authGroupV1.Put("/password/forgot", userHandler.ResetPassword)
	authGroupV1.Put("/password/reset", userHandler.ConfirmPassowrd)
}
