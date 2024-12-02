package types

import "github.com/gofiber/fiber/v2"

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func SendResponse(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return ctx.Status(statusCode).JSON(Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
