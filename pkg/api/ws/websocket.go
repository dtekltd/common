package wsApi

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/pkg/ws"
	"github.com/gofiber/fiber/v2"
)

func sendMessage(ctx *fiber.Ctx) error {
	msg := &ws.Message{}
	if err := ctx.BodyParser(msg); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if err := ws.SendMessage(msg); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	return api.SuccessResp(ctx, true)
}
