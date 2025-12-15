package cmsApi

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func getAnnounces(ctx *fiber.Ctx, postType string) error {
	var rows []*cms.Post

	if err := database.DB.Model(&cms.Post{}).
		Where("(type=?) AND (published_at>0)", postType).
		Select("id", "type", "name", "content", "meta", "created_at", "published_at", "created_by").
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Order("created_at DESC").
		Limit(ctx.QueryInt("limit", 3)).
		Find(&rows).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return api.SuccessResp(ctx, rows)
}
