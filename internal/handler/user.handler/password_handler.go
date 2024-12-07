package userhandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.SendEmail)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	var ErrorMessages []types.ErrorField
	if err := request.Validate(); err != nil {
		for _, errs := range err.(validator.ValidationErrors) {
			var errMsg types.ErrorField
			errMsg.FailedField = errs.Field()
			errMsg.Tag = errs.Tag()
			errMsg.Value = errs.Value()

			ErrorMessages = append(ErrorMessages, errMsg)
		}
		logrus.WithField("validate body", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", ErrorMessages)
	}

	if err := h.service.ResetPassword(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("Reset Password", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success send token", nil)
}

func (h *userHandler) ConfirmPassowrd(ctx *fiber.Ctx) error {
	request := new(model.ResetPassword)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	var ErrorMessages []types.ErrorField
	if err := request.Validate(); err != nil {
		for _, errs := range err.(validator.ValidationErrors) {
			var errMsg types.ErrorField
			errMsg.FailedField = errs.Field()
			errMsg.Tag = errs.Tag()
			errMsg.Value = errs.Value()

			ErrorMessages = append(ErrorMessages, errMsg)
		}
		logrus.WithField("validate body", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", ErrorMessages)
	}

	if err := h.service.ConfirmPassword(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("Confirm Password", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success updated password", nil)
}
