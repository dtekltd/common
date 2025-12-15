package api

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResp(ctx *fiber.Ctx, data any, meta ...ApiResponseMeta) error {
	resp := ApiResponse{
		Success: true,
		Data:    data,
	}
	if len(meta) > 0 {
		resp.Meta = &meta[0]
	}
	return ctx.Status(fiber.StatusOK).JSON(&resp)
}

func ErrorResp(ctx *fiber.Ctx, err ApiError, meta ...ApiResponseMeta) error {
	resp := ApiResponse{
		Success: false,
		Error:   &err,
	}
	if len(meta) > 0 {
		resp.Meta = &meta[0]
	}
	code := fiber.StatusBadRequest
	if err.Code != 0 {
		code = err.Code
	}
	return ctx.Status(code).JSON(&resp)
}

func ErrorCodeResp(ctx *fiber.Ctx, code int, message ...string) error {
	msg := "API Error"
	if len(message) > 0 {
		msg = message[0]
	}
	return ErrorResp(ctx, ApiError{
		Code:    code,
		Message: msg,
	})
}

func ErrorNotFoundResp(ctx *fiber.Ctx, message ...string) error {
	return ErrorCodeResp(ctx, fiber.StatusNotFound, message...)
}

func ErrorUnauthorizedResp(ctx *fiber.Ctx, message ...string) error {
	return ErrorCodeResp(ctx, fiber.StatusUnauthorized, message...)
}

func ErrorBadRequestResp(ctx *fiber.Ctx, message ...string) error {
	return ErrorCodeResp(ctx, fiber.StatusBadRequest, message...)
}

func ErrorInternalServerErrorResp(ctx *fiber.Ctx, message ...string) error {
	return ErrorCodeResp(ctx, fiber.StatusInternalServerError, message...)
}
