package cmsApi

import (
	"slices"
	"time"

	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/strcase"
	"github.com/dtekltd/common/types"
	"github.com/dtekltd/common/utils"
	"github.com/gofiber/fiber/v2"
)

func withPostType(handler func(*fiber.Ctx, string) error, postType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handler(c, postType)
	}
}

func getPost(ctx *fiber.Ctx, postType string) error {
	post := &cms.Post{}
	if err := database.DB.Preload("Categories").
		Take(post, ctx.Params("id")).Error; err != nil {
		return api.ErrorNotFoundResp(ctx, "Post not found!")
	}
	if post.Type != postType {
		return api.ErrorNotFoundResp(ctx, "Post not found!")
	}
	return api.SuccessResp(ctx, post)
}

func getPosts(ctx *fiber.Ctx, postType string) error {
	type PostQueryFilter struct {
		Target string `json:"target"`
		Status int64  `json:"status"`
	}

	// type PostQueryParams struct {
	// 	Page    int             `json:"page"`
	// 	PerPage int             `json:"perPage"`
	// 	Filter  PostQueryFilter `json:"filter"`
	// }

	filter := &PostQueryFilter{}
	if err := utils.QueryStruct(ctx, filter, "filter"); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	var count int64
	var rows []*cms.Post

	query := database.DB.Model(&cms.Post{}).Where("type=?", postType)
	_ = query.Count(&count)

	if filter.Target == "" {
		// select only required fields
		query = query.Select("id", "type", "name", "created_at", "published_at", "ordering")
	}

	ordering := &types.Params{}
	pagination := &api.Pagination{
		Total: int(count),
	}
	if err := query.
		Scopes(database.Ordering(ctx, ordering)).
		Scopes(database.Paginate(ctx, pagination)).
		Find(&rows).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return api.SuccessResp(ctx, rows, api.ApiResponseMeta{
		Ordering:   ordering,
		Pagination: pagination,
	})
}

func savePost(ctx *fiber.Ctx, postType string) error {
	post := &cms.Post{}
	if err := ctx.BodyParser(post); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	if post.Path == "" {
		post.Path = strcase.KebabCase(post.Name)
	}

	tx := database.DB.Begin()
	if post.ID == 0 {
		post.CreatedBy = ctx.Locals("uiID").(uint64)
		post.CreatedAt = 0
		post.PublishedAt = 0
	} else {
		if post.Type == "post" {
			// handle deleted categories
			catIDs := []uint64{}
			model := &cms.PostCategory{}
			database.DB.Model(model).
				Select("category_id").
				Where("post_id=?", post.ID).
				Find(&catIDs)

			if len(catIDs) > 0 {
				for _, cat := range post.Categories {
					idx := slices.Index(catIDs, cat.ID)
					if idx != -1 {
						catIDs = slices.Delete(catIDs, idx, idx+1)
					}
				}
				if len(catIDs) > 0 {
					if err := tx.Delete(model, "category_id IN (?)", catIDs).Error; err != nil {
						return api.ErrorInternalServerErrorResp(ctx, err.Error())
					}
				}
			}
		}
	}

	if err := tx.Save(post).Error; err != nil {
		tx.Rollback()
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		tx.Commit()
		return api.SuccessResp(ctx, post)
	}
}

func publishPost(ctx *fiber.Ctx, postType string) error {
	post := &cms.Post{}
	if err := database.DB.Take(post, ctx.Params("id")).Error; err != nil {
		return api.ErrorNotFoundResp(ctx, "Post not found!")
	}
	if post.Type != postType {
		return api.ErrorNotFoundResp(ctx, "Post not found!")
	}
	if post.PublishedAt == 0 {
		post.PublishedAt = uint64(time.Now().Unix())
	} else {
		post.PublishedAt = 0
	}
	if err := database.DB.Save(post).Error; err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	} else {
		return api.SuccessResp(ctx, post)
	}
}
