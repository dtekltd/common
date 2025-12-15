package mailApi

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/mail"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
)

func getInstances(ctx *fiber.Ctx) error {
	if id, err := ctx.ParamsInt("id"); err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	} else {
		var count int64
		var rows []*mail.Instance

		query := mail.DB().Model(&mail.Instance{}).
			Where("message_id=?", id)
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
			Find(&rows).Error
		if err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		return api.SuccessResp(ctx, rows, api.ApiResponseMeta{
			Ordering:   ordering,
			Pagination: pagination,
		})
	}
}

func deleteInstances(ctx *fiber.Ctx) error {
	if id, err := ctx.ParamsInt("id"); err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	} else {
		where := []any{
			"(message_id=?)",
			id,
		}
		if target := ctx.Query("target"); target != "" {
			if target == "sent" {
				where[0] = where[0].(string) + " AND (sent_at>0)"
			}
		}
		if err := mail.DB().Delete(&mail.Instance{}, where...).Error; err != nil {
			return api.ErrorBadRequestResp(ctx, err.Error())
		}
		return api.SuccessResp(ctx, true)
	}
}
