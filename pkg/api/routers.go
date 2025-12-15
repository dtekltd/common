package api

import (
	"github.com/dtekltd/common/jwt"
	accountsApi "github.com/dtekltd/common/pkg/api/accounts"
	authApi "github.com/dtekltd/common/pkg/api/auth"
	cmsApi "github.com/dtekltd/common/pkg/api/cms"
	mailApi "github.com/dtekltd/common/pkg/api/mail"
	mediaApi "github.com/dtekltd/common/pkg/api/media"
	siteApi "github.com/dtekltd/common/pkg/api/site"
	wsApi "github.com/dtekltd/common/pkg/api/ws"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router, keyManager *jwt.KeyManager) {
	// router := app.Group("/api", auth.AuthMiddleware(keyManager, "/api"))

	siteApi.RegisterHandlers(router, keyManager)
	authApi.RegisterHandlers(router, keyManager)
	accountsApi.RegisterHandlers(router, keyManager)
	mailApi.RegisterHandlers(router)
	mediaApi.RegisterHandlers(router)
	cmsApi.RegisterHandlers(router)
	wsApi.RegisterHandlers(router)
}
