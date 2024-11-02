package middleware

import (
	"kisara/src/response"
	"kisara/src/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorization := c.GetRespHeader("Authorization")

		if authorization == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusUnauthorized,
				Name:       "Unauthorized",
				Message:    "No token provided",
			})
		}

		if strings.Split(authorization, " ")[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusUnauthorized,
				Name:       "Unauthorized",
				Message:    "Invalid token",
			})
		}

		authorization = strings.Split(authorization, " ")[1]

		// Check if the token is valid
		token, err := utils.ValidateToken(authorization)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusUnauthorized,
				Name:       "Unauthorized",
				Message:    "Invalid token",
			})
		}

		c.Locals("token", token)

		return c.Next()
	}
}
