package userhandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.SendEmail)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := request.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := h.service.ResetPassword(ctx.Context(), *request); err != nil {
		logrus.WithField("Reset Password", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success send token", nil)
}

func (h *userHandler) ConfirmPassowrd(ctx *fiber.Ctx) error {
	request := new(model.ResetPassword)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := request.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := h.service.ConfirmPassword(ctx.Context(), *request); err != nil {
		logrus.WithField("Confirm Password", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success updated password", nil)
}
