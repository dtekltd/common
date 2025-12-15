package mailApi

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/mail"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
)

func getMessage(ctx *fiber.Ctx) error {
	if id, err := ctx.ParamsInt("id"); err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	} else {
		msg := &mail.Message{}
		if err := mail.DB().Take(msg, id).Error; err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}
		return api.SuccessResp(ctx, msg)
	}
}

func getMessages(ctx *fiber.Ctx) error {
	var count int64
	var rows []*mail.Message

	query := mail.DB().Model(&mail.Message{})
	_ = query.Count(&count)

	ordering := &types.Params{}
	if ctx.Query("orderBy") == "" {
		(*ordering)["created_at"] = "DESC" // default ordering
	}
	pagination := &api.Pagination{
		Total: int(count),
	}
	err := query.
		Scopes(database.Ordering(ctx, ordering)).
		Scopes(database.Paginate(ctx, pagination)).
		Select("id", "type", "subject", "frequency", "created_at", "processed_at").
		Find(&rows).Error
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return api.SuccessResp(ctx, rows, api.ApiResponseMeta{
		Ordering:   ordering,
		Pagination: pagination,
	})
}

func saveMessage(ctx *fiber.Ctx) error {
	msg := &mail.Message{}
	if err := ctx.BodyParser(msg); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if msg.Subject == "" || msg.Body == "" {
		return api.ErrorBadRequestResp(ctx, "Both \"subject\" and \"body\" are required.")
	}
	// validate template if any
	if err := msg.ParseBodyTmplt(); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}
	if err := msg.Save(); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, msg)
	}
}

func sendMessage(ctx *fiber.Ctx) error {
	if id, err := ctx.ParamsInt("id"); err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	} else {
		msg := &mail.Message{}
		if err := mail.DB().Take(msg, id).Error; err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		params := types.Params{}
		if force := ctx.Query("force"); force == "true" {
			params["__force"] = true
		}
		if immediate := ctx.Query("immediate"); immediate == "true" {
			params["__immediate"] = true
		}
		if err := mail.SendMessage(msg, params); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}
		return api.SuccessResp(ctx, msg)
	}
}

func deleteMessage(ctx *fiber.Ctx) error {
	if id, err := ctx.ParamsInt("id"); err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	} else {
		if err := mail.DB().Delete(&mail.Message{}, uint64(id)).Error; err != nil {
			return api.ErrorBadRequestResp(ctx, err.Error())
		}
		return api.SuccessResp(ctx, true)
	}
}
