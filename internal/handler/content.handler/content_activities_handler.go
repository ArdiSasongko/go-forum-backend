package contenthandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *contentHandler) UpdateUserActivitiesContent(ctx *fiber.Ctx) error {
	request := new(model.LikedModel)
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
	request.UserID = ctx.Locals("user_id").(int32)
	request.Username = ctx.Locals("username").(string)
	contentID, _ := ctx.ParamsInt("content_id")
	request.ContentID = int32(contentID)

	if err := h.service.LikedDislikeContent(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("liked content", err.Error()).Error("failed to liked content")
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success liked", nil)
}
