package userhandler

import (
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) GetProfile(ctx *fiber.Ctx) error {
	username := ctx.Locals("username").(string)
	result, err := h.service.GetProfile(ctx.Context(), username)
	if err == sql.ErrNoRows {
		logrus.WithField("get profile", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}
	if err != nil {
		logrus.WithField("get profile", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success get profile", result)
}
