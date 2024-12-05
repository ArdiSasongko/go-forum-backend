package userhandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) ValidateUser(ctx *fiber.Ctx) error {
	request := new(model.ValidateToken)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := request.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	payload := model.ValidatePayload{
		Token:    request.Token,
		Username: ctx.Locals("username").(string),
	}

	if err := h.service.ValidateEmail(ctx.Context(), queries, payload); err != nil {
		logrus.WithField("validate email", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success validate email", nil)
}

func (h *userHandler) ResendEmail(ctx *fiber.Ctx) error {
	payload := model.ValidatePayload{
		Username: ctx.Locals("username").(string),
	}

	if err := h.service.ResendEmail(ctx.Context(), queries, payload); err != nil {
		logrus.WithField("resend email", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success resend email", nil)
}
