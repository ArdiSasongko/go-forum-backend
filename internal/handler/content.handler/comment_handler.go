package contenthandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *contentHandler) InsertComment(ctx *fiber.Ctx) error {
	request := new(model.CommentModel)
	contentID, _ := ctx.ParamsInt("content_id")
	username := ctx.Locals("username").(string)
	user_id := ctx.Locals("user_id").(int32)

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

	request.UserID = user_id
	request.ContentID = int32(contentID)
	request.Username = username
	if err := h.service.InsertComment(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("insert comment", err.Error()).Error("failed to insert comment")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusCreated, "success insert comment", nil)
}

func (h *contentHandler) DeleteComment(ctx *fiber.Ctx) error {
	contentID, _ := ctx.ParamsInt("content_id")
	user_id := ctx.Locals("user_id").(int32)

	if err := h.service.DeleteComment(ctx.Context(), queries, user_id, int32(contentID)); err != nil {
		logrus.WithField("delete comment", err.Error()).Error("failed to failed comment")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success delete comment", nil)
}
