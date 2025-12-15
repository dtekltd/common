package accountsApi

import (
	"strings"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/auth"
	"github.com/dtekltd/common/pkg/users"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func filter(ctx *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		filter := ctx.Query("filter")
		if filter != "" {
			search := "%" + strings.Replace(filter, "%", "\\%", -1) + "%"
			return db.Where("(name LIKE ?) OR (email like ?)", search, search)
		}
		return db
	}
}

func getAccounts(ctx *fiber.Ctx) error {
	var count int64
	var rows []*users.Account

	query := database.DB.Model(&users.Account{}).
		Scopes(filter(ctx))
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
		Select("id", "username", "name", "email", "created_at").
		Find(&rows).Error
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return api.SuccessResp(ctx, rows, api.ApiResponseMeta{
		Ordering:   ordering,
		Pagination: pagination,
	})
}

func createAccount(ctx *fiber.Ctx) error {
	req := &auth.RegisterReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if req.Email == "" || req.Password == "" {
		return api.ErrorBadRequestResp(ctx, "Both \"email\" and \"password\" are required.")
	}
	if acc, err := auth.Register(req); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, acc)
	}
}

func deleteAccount(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return api.ErrorBadRequestResp(ctx, "Invalid ID")
	}
	if err := database.DB.Delete(&users.Account{
		ID: uint64(id),
	}).Error; err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	return api.SuccessResp(ctx, true)
}
