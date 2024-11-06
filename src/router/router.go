package router

import (
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(response.GeneralResponse{
			StatusCode: fiber.StatusOK,
			Name:       "Sleep. Zzz...",
			Message:    "Welcome to the Kisara API",
		})
	})

	SetupAuthRoutes(app, db)
	SetupUserRoutes(app, db)
	SetupMessageRoutes(app, db)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(response.NotFoundResponse{
			GeneralResponse: response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    "Route not found",
			},
			Path: c.Path(),
		})
	})
}
