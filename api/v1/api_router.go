package v1

import (
	contenthandler "github.com/ArdiSasongko/go-forum-backend/internal/handler/content.handler"
	userhandler "github.com/ArdiSasongko/go-forum-backend/internal/handler/user.handler"
	contentservice "github.com/ArdiSasongko/go-forum-backend/internal/service/content.service"
	userservice "github.com/ArdiSasongko/go-forum-backend/internal/service/user.service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type ApiRouter struct {
	userService    userservice.UserService
	contentService contentservice.ContentService
}

func NewApiRouter(
	userService userservice.UserService,
	contentService contentservice.ContentService,
) *ApiRouter {
	return &ApiRouter{userService: userService, contentService: contentService}
}

func (h *ApiRouter) InstallRouter(app *fiber.App) {
	h.setupAuthRouter(app)
	h.setupUserRouter(app)
	h.setupContentRouter(app)
}

func (h *ApiRouter) setupAuthRouter(app *fiber.App) {
	userHandler := userhandler.NewUserHandler(h.userService)
	authGroup := app.Group("/auth")
	authGroupV1 := authGroup.Group("/v1")

	authGroupV1.Post("/register", userHandler.Register)
	authGroupV1.Post("/login", userHandler.Login)
	authGroupV1.Put("/refresh-token", middleware.MiddlewareRefreshToken, userHandler.RefreshToken)
	authGroupV1.Put("/validate/email", middleware.MiddlewareAuthValidate, userHandler.ValidateUser)
	authGroupV1.Get("/validate/resend", middleware.MiddlewareAuthValidate, userHandler.ResendEmail)
	authGroupV1.Put("/password/forgot", userHandler.ResetPassword)
	authGroupV1.Put("/password/reset", userHandler.ConfirmPassowrd)
}

func (h *ApiRouter) setupUserRouter(app *fiber.App) {
	userHandler := userhandler.NewUserHandler(h.userService)
	userGroup := app.Group("/user")
	userGroupV1 := userGroup.Group("/v1")

	userGroupV1.Get("/profile", middleware.MiddlewareAuthValidate, userHandler.GetProfile)
	userGroupV1.Put("/profile", middleware.MiddlewareAuthValidate, userHandler.UpdateUser)
	userGroupV1.Put("/profile/image", middleware.MiddlewareAuthValidate, userHandler.UpdateProfile)
}

func (h *ApiRouter) setupContentRouter(app *fiber.App) {
	contentHandler := contenthandler.NewContentHandler(h.contentService)
	contentGroup := app.Group("/content")
	contentGroupV1 := contentGroup.Group("/v1")

	contentGroupV1.Post("/", middleware.MiddlewareAuthValidate, contentHandler.CreateContent)
	contentGroupV1.Get("/", contentHandler.GetContents)
	contentGroupV1.Get("/:content_id", contentHandler.GetContent)
	contentGroupV1.Get("/page/:page", contentHandler.GetContents)
}
