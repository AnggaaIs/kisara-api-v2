package router

import (
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Welcome to the Kisara API",
			nil,
		))
	})

	SetupAuthRoutes(app, db)
	SetupUserRoutes(app, db)
	SetupMessageRoutes(app, db)

	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(
			fiber.StatusNotFound,
			"Not Found",
			"Route not found",
			nil,
		))
	})
}
