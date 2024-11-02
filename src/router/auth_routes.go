package router

import (
	"kisara/src/controller"
	"kisara/src/middleware"
	"kisara/src/models/validation"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	auth := app.Group("/auth")

	auth.Get("/google/url", controller.HandleGoogleURL(db))
	auth.Post("/google/callback", controller.HandleGoogleCallback(db), middleware.ValidateSchemas(nil, validation.AuthGoogleCallbackBody{}))
}
