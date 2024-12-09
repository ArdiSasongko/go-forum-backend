package contenthandler

import (
	"database/sql"

	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *contentHandler) CreateContent(ctx *fiber.Ctx) error {
	request := new(model.ContentModel)

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
	if err == nil && len(files.File["file"]) > 0 {
		request.Files = files.File["file"]
	}

	request.UserID = ctx.Locals("user_id").(int32)
	request.Username = ctx.Locals("username").(string)

	if err := h.service.InsertContent(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("create content", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusCreated, "success created content", nil)
}

func (h *contentHandler) GetContents(ctx *fiber.Ctx) error {
	page, err := ctx.ParamsInt("page")

	if err != nil || page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	contents, err := h.service.GetContents(ctx.Context(), queries, int32(limit), int32(offset))
	if err != nil {
		logrus.WithField("get contents", "BAD REQUEST").Error("failed to get contents")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	response := fiber.Map{
		"contents": contents,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
		},
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success get contents", response)
}

func (h *contentHandler) GetContent(ctx *fiber.Ctx) error {
	contentID, _ := ctx.ParamsInt("content_id")
	userID := ctx.Locals("user_id").(int32)
	page := 1
	limit := 10
	offset := (page - 1) * limit
	content, err := h.service.GetContent(ctx.Context(), queries, int32(contentID), userID, int32(offset), int32(limit))
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithField("get contents", "BAD REQUEST").Error("failed to get content")
			return types.SendResponse(ctx, fiber.StatusNotFound, "BAD REQUEST", sql.ErrNoRows)
		}
		logrus.WithField("get contents", "BAD REQUEST").Error("failed to get content")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success get content", content)
}

func (h *contentHandler) UpdateContent(ctx *fiber.Ctx) error {
	request := new(model.UpdateContent)
	contentID, _ := ctx.ParamsInt("content_id")

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

	username := ctx.Locals("username").(string)
	userID := ctx.Locals("user_id").(int32)

	request.UpdatedBy = username
	if err := h.service.UpdateContent(ctx.Context(), queries, int32(contentID), userID, *request); err != nil {
		logrus.WithField("update content", err.Error()).Error("failed to update content")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success update content", nil)
}

func (h *contentHandler) DeleteContent(ctx *fiber.Ctx) error {
	contentID, _ := ctx.ParamsInt("content_id")

	if err := h.service.DeleteContent(ctx.Context(), queries, int32(contentID)); err != nil {
		logrus.WithField("delete content", err.Error()).Error("failed to delete content")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success delete content", nil)
}
