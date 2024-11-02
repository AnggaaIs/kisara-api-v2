package controller

import (
	"kisara/src/models"
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func HandleGetUser(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		email := claims["email"].(string)

		var user models.User
		result := db.Where("email = ?", email).Select("name", "link_id", "profile_url").First(&user)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    "User not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.DataResponse{
			GeneralResponse: response.GeneralResponse{
				StatusCode: fiber.StatusOK,
				Name:       "Success",
				Message:    "User found",
			},
			Data: user,
		})
	}
}
