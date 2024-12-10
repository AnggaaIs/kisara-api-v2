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
	message.Get("/:link_id", controller.HandleMessageGet(db), middleware.ValidateSchemas(validation.MessageBodyGet{}, nil))
	message.Delete("/:link_id/:message_id", controller.HandleDeleteMessage(db), middleware.AuthMiddleware(db))

	message.Get("/:link_id/:message_id", controller.HandleReplyMessageGet(db), middleware.ValidateSchemas(validation.MessageBodyGet{}, nil))
	message.Post("/:link_id/:message_id", controller.HandleReplyMessagePost(db), middleware.AuthMiddleware(db), middleware.ValidateSchemas(nil, validation.MessageBodyPost{}))
	message.Delete("/:link_id/:message_id/:reply_id", controller.HandleDeleteReplyMessage(db), middleware.AuthMiddleware(db))

	message.Put("/:link_id/:message_id/like", controller.HandleLikeMessage(db), middleware.AuthMiddleware(db))
}
