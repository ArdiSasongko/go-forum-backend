package handler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/ArdiSasongko/go-forum-backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Register(ctx *fiber.Ctx) error {
	user := new(model.UserModel)

	if err := ctx.BodyParser(user); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := user.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := h.service.CreateUser(ctx.Context(), *user); err != nil {
		logrus.WithField("create user", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusBadRequest, "success register", nil)
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	user := new(model.LoginRequest)

	if err := ctx.BodyParser(user); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := user.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	token, err := h.service.LoginUser(ctx.Context(), *user)
	if err != nil {
		logrus.WithField("login user", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success login", token)
}
