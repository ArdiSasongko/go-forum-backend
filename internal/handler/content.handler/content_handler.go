package contenthandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *contentHandler) CreateContent(ctx *fiber.Ctx) error {
	request := new(model.ContentModel)

	if err := ctx.BodyParser(request); err != nil {
		logrus.WithField("parsing body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if err := request.Validate(); err != nil {
		logrus.WithField("validate body", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	files, err := ctx.MultipartForm()
	if err == nil && len(files.File["file"]) > 0 {
		request.Files = files.File["file"]
	}

	request.UserID = ctx.Locals("user_id").(int32)
	request.Username = ctx.Locals("username").(string)

	if err := h.service.InsertContent(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("create content", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
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
		logrus.WithField("get contents", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
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
