package wsApi

import (
	"github.com/dtekltd/common/pkg/ws"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router) {
	router = router.Group("/ws")

	router.Get("/:id/connect", ws.NewHandler())
	router.Post("/message", sendMessage)

	// run WS hub in a new routine
	go ws.RunHub()
}
