package mailApi

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router) {
	mail := router.Group("/mail")

	messages := mail.Group("/messages")
	messages.Get("/:id", getMessage)
	messages.Get("", getMessages)
	messages.Post("", saveMessage)
	messages.Put("/:id", sendMessage)
	messages.Delete("/:id", deleteMessage)

	instances := mail.Group("/instances")
	instances.Get("/:id", getInstances)
	instances.Delete("/:id", deleteInstances)
}
