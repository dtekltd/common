package cmsApi

import (
	"time"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/strcase"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
)

func getCategory(ctx *fiber.Ctx, catType string) error {
	cat := &cms.Category{}
	if err := database.DB.Take(cat, ctx.Params("id")).Error; err != nil {
		return api.ErrorNotFoundResp(ctx, "Category not found!")
	}
	if cat.Type != catType {
		return api.ErrorNotFoundResp(ctx, "Category not found!")
	}
	return api.SuccessResp(ctx, cat)
}

func getCategories(ctx *fiber.Ctx, catType string) error {
	var count int64
	var rows []*cms.Category

	query := database.DB.Model(&cms.Category{}).
		Where("type=?", catType)
	_ = query.Count(&count)

	ordering := &types.Params{}
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

func saveCategory(ctx *fiber.Ctx, catType string) error {
	cat := &cms.Category{}
	if err := ctx.BodyParser(cat); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if cat.Alias == "" {
		cat.Alias = strcase.KebabCase(cat.Name)
	}
	if err := database.DB.Save(cat).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, cat)
	}
}

func publishCategory(ctx *fiber.Ctx, catType string) error {
	cat := &cms.Category{}
	if err := database.DB.Take(cat, ctx.Params("id")).Error; err != nil {
		return api.ErrorNotFoundResp(ctx, "Category not found!")
	}
	if cat.Type != catType {
		return api.ErrorNotFoundResp(ctx, "Category not found!")
	}
	if cat.PublishedAt == 0 {
		cat.PublishedAt = uint64(time.Now().Unix())
	} else {
		cat.PublishedAt = uint64(0)
	}
	if err := database.DB.Save(cat).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, cat)
	}
}
