package router

import (
	"kisara/src/controller"
	"kisara/src/middleware"
	"kisara/src/models/validation"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	auth := app.Group("/auth", middleware.RateLimitMiddleware(middleware.RateLimitConfig{
		Max:      15,
		Duration: 60,
	}))

	auth.Get("/google/url", controller.HandleGoogleURL(db))
	auth.Post("/google/callback", controller.HandleGoogleCallback(db), middleware.ValidateSchemas(nil, validation.AuthGoogleCallbackBody{}))
}
