package cms

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/system"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBlogPosts(ctx *fiber.Ctx, postTypes ...string) types.Params {
	posts := []types.Params{}

	var count int64
	var rows []*Post

	postType := "post"
	if len(postTypes) > 0 {
		postType = postTypes[0]
	}
	query := database.DB.Model(&Post{}).
		Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "alias")
		}).
		Select("id", "name", "path", "meta", "created_at", "updated_at", "published_at").
		Where("(type=?) AND (published_at > 0)", postType).
		Order("created_at DESC")
	_ = query.Count(&count)

	pagination := &api.Pagination{
		Total: int(count),
	}
	if err := query.
		Scopes(database.Paginate(ctx, pagination)).
		Find(&rows).Error; err != nil {
		system.Logger.Error("cms::GetBlogPosts error", err.Error())
	}

	for _, row := range rows {
		posts = append(posts, MapPost(row))
	}

	return types.Params{
		"Count": count,
		"Posts": posts,
	}
}
