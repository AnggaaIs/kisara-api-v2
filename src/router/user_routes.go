package router

import (
	"kisara/src/controller"
	"kisara/src/middleware"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	user := app.Group("/user")

	user.Get("/", controller.HandleGetUser(db), middleware.AuthMiddleware(db))
}
