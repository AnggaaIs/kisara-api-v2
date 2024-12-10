package controller

import (
	"kisara/src/models"
	"kisara/src/models/validation"
	"kisara/src/response"
	"reflect"

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

func HandleUpdateUser(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		email := claims["email"].(string)
		data_update := c.Locals("requestBody").(*validation.UserUpdateBody)

		var user models.User
		result := db.Model(&models.User{}).
			Where("email = ?", email).
			First(&user)

		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				"User not found",
				result.Error,
			))
		}

		oldUser := user

		userValue := reflect.ValueOf(&user).Elem()
		dataUpdateValue := reflect.ValueOf(data_update).Elem()

		updated := false

		for i := 0; i < dataUpdateValue.NumField(); i++ {
			fieldName := dataUpdateValue.Type().Field(i).Name
			fieldValue := dataUpdateValue.Field(i).Interface()

			if fieldValue != "" {
				userField := userValue.FieldByName(fieldName)
				if userField.IsValid() && userField.CanSet() && userField.String() != fieldValue {
					userField.Set(reflect.ValueOf(fieldValue))
					updated = true
				}
			}
		}

		if !updated {
			return c.Status(fiber.StatusOK).JSON(response.Success(
				fiber.StatusOK,
				"No Update",
				"No changes made to the user",
				nil,
			))
		}

		result = db.Save(&user)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to update user",
				result.Error,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"User updated",
			map[string]interface{}{
				"old_data": oldUser,
				"new_data": user,
			},
		))
	}
}
