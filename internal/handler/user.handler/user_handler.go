package userhandler

import (
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) GetProfile(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	result, err := h.service.GetProfile(ctx.Context(), queries, email)
	if err == sql.ErrNoRows {
		logrus.WithField("get profile", "NOT FOUND").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusNotFound, "NOT FOUND", sql.ErrNoRows)
	}
	if err != nil {
		logrus.WithField("get profile", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success get profile", result)
}

func (h *userHandler) UpdateProfile(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	request := new(model.UpdateProfile)

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

	files, err := ctx.MultipartForm()
	if err == nil {
		if len(files.File["file"]) > 1 {
			logrus.WithField("update image profile", "BAD REQUEST").Error("only one can upload to updated")
			return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", "only one can upload to updated")
		}

		if len(files.File["file"]) == 1 {
			request.Files = files.File["file"]
		}
	}

	request.Email = email

	if err := h.service.UpdateProfile(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("update image profile", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusCreated, "success update profile image", nil)
}

func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	request := new(model.UpdateUser)
	email := ctx.Locals("email").(string)

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

	if err := h.service.UpdateUser(ctx.Context(), queries, *request, email); err != nil {
		logrus.WithField("update user", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success update profile", nil)
}
