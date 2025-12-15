package accountsApi

import (
	"github.com/dtekltd/common/jwt"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router, keyManager *jwt.KeyManager) {
	// profile
	profile := router.Group("/profile")
	profile.Put("", saveAccount)
	profile.Post("/avatar", uploadAvatar)

	// accounts
	accounts := router.Group("/accounts")

	accounts.Get("", getAccounts)
	accounts.Post("", createAccount)
	accounts.Put("/:id?", saveAccount)
	accounts.Delete("/:id", deleteAccount)
}
