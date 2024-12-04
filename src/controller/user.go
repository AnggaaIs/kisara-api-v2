package controller

import (
	"kisara/src/models"
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserResponse struct {
	Name       string `json:"name"`
	LinkID     string `json:"link_id"`
	ProfileURL string `json:"profile_url"`
}

func HandleGetUser(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		email := claims["email"].(string)

		var userResponse UserResponse
		result := db.Model(&models.User{}).
			Where("email = ?", email).
			Select("name, link_id, profile_url").
			Scan(&userResponse)

		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				"User not found",
				result.Error,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"User found",
			userResponse,
		))
	}
}
