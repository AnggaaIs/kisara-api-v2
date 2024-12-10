package middleware

import (
	"kisara/src/response"
	"kisara/src/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		authorization := c.Get("Authorization")

		if authorization == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"No token provided",
				nil,
			))
		}

		if strings.Split(authorization, " ")[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"Invalid token",
				nil,
			))
		}

		authorization = strings.Split(authorization, " ")[1]

		// Check if the token is valid
		token, err := utils.ValidateToken(authorization)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"Invalid token",
				err,
			))
		}

		claims := token.Claims.(jwt.MapClaims)

		timeToken := claims["time"].(float64)

		// Check if the token is expired (7 days)
		if timeToken+604800 < float64(time.Now().Unix()) {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"Token expired",
				nil,
			))
		}

		c.Locals("token", token)

		return c.Next()
	}
}
