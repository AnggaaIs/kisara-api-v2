package controller

import (
	"kisara/src/models"
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func HandleStats(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var totalUsers int64
		if err := db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to retrieve total users",
				err,
			))
		}

		var totalComments int64
		if err := db.Model(&models.Comment{}).Count(&totalComments).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to retrieve total comments",
				err,
			))
		}

		stats := map[string]interface{}{
			"total_users":    totalUsers,
			"total_comments": totalComments,
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"API Statistics Retrieved Successfully",
			stats,
		))
	}
}
