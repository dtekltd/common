package siteApi

import (
	"github.com/dtekltd/common/api"
	"github.com/dtekltd/common/database"
	"github.com/dtekltd/common/jwt"
	"github.com/dtekltd/common/pkg/cms"
	"github.com/dtekltd/common/pkg/site"
	"github.com/dtekltd/common/types"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router, keyManager *jwt.KeyManager) {
	router = router.Group("/site")

	router.Get("/", siteSettings)
	router.Post("/", updateSiteSettings)

	router.Get("/agreement", agreementArticle)
}

func siteSettings(ctx *fiber.Ctx) error {
	return api.SuccessResp(ctx, site.Settings())
}

func updateSiteSettings(ctx *fiber.Ctx) error {
	params := site.Settings()
	// need to reset custom_params
	// or, it will be merged with the old params
	params.CustomParams = &types.Params{}
	if err := ctx.BodyParser(params); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}
	err := site.UpdateSettings(params)
	return api.SuccessResp(ctx, err == nil)
}

func agreementArticle(ctx *fiber.Ctx) error {
	if postID := site.Settings().CustomParams.GetInt("agreementArticle"); postID != 0 {
		post := &cms.Post{}
		if err := database.DB.Take(post, postID).Error; err == nil {
			return api.SuccessResp(ctx, post)
		}
	}
	return api.SuccessResp(ctx, "No agreement")
}
