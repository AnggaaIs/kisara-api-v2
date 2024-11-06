package router

import (
	"kisara/src/controller"
	"kisara/src/middleware"
	"kisara/src/models/validation"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupMessageRoutes(app *fiber.App, db *gorm.DB) {
	message := app.Group("/message", middleware.RateLimitMiddleware(middleware.RateLimitConfig{
		Max:      30,
		Duration: 60,
	}))

	message.Post("/:link_id", controller.HandleMessagePost(db), middleware.ValidateSchemas(nil, validation.MessageBodyPost{}))
	message.Get("/:link_id", controller.HandleMessageGet(db))
	message.Delete("/:link_id/:message_id", controller.HandleDeleteMessage(db), middleware.AuthMiddleware(db))
}
