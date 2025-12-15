package mediaApi

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router) {
	media := router.Group("/media")

	media.Get("", getMedia)
	media.Put("", createFolder)
	media.Post("", uploadMedia)
	media.Delete("", deleteMedia)
}
