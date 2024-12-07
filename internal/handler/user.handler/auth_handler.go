package userhandler

import (
	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func (h *userHandler) Register(ctx *fiber.Ctx) error {
	request := new(model.UserModel)

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

	if err := h.service.CreateUser(ctx.Context(), queries, *request); err != nil {
		logrus.WithField("create user", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success register", nil)
}

func (h *userHandler) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

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

	token, err := h.service.LoginUser(ctx.Context(), queries, *request)
	if err != nil {
		logrus.WithField("login user", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success login", token)
}

func (h *userHandler) RefreshToken(ctx *fiber.Ctx) error {
	request := new(model.RefreshToken)
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

	payload := model.PayloadToken{
		UserID:   ctx.Locals("user_id").(int32),
		Username: ctx.Locals("username").(string),
		Email:    ctx.Locals("email").(string),
		Role:     ctx.Locals("role").(string),
	}

	newToken, err := h.service.RefreshToken(ctx.Context(), queries, payload, *request)
	if err != nil {
		logrus.WithField("get token", "BAD REQUEST").Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusBadRequest, "BAD REQUEST", err.Error())
	}

	return types.SendResponse(ctx, fiber.StatusOK, "success get token", newToken)
}
