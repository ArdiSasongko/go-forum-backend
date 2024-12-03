package router

import (
	"github.com/ArdiSasongko/go-forum-backend/internal/handler"
	"github.com/ArdiSasongko/go-forum-backend/internal/service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type ApiRouter struct {
	userService service.UserService
}

func NewApiRouter(userService service.UserService) *ApiRouter {
	return &ApiRouter{userService: userService}
}

func (h ApiRouter) InstallRouter(app *fiber.App) {
	authGroup := app.Group("/user")
	authGroupV1 := authGroup.Group("/v1")

	userHandler := handler.NewUserHandler(h.userService)
	authGroupV1.Post("/register", userHandler.Register)
	authGroupV1.Post("/login", userHandler.Login)
	authGroupV1.Put("/refresh-token", middleware.MiddlewareRefreshToken, userHandler.RefreshToken)
}
