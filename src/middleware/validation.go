// middleware/validator.go
package middleware

import (
	"kisara/src/response"
	"kisara/src/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateSchemas(querySchema interface{}, bodySchema interface{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		validate := validator.New()

		if querySchema != nil {
			queryInstance := utils.CreateNew(querySchema)

			if err := c.Bind().Query(queryInstance); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(response.Error(
					fiber.StatusBadRequest,
					"Bad Request",
					"Invalid query parameters",
					nil,
				))
			}

			if err := validate.Struct(queryInstance); err != nil {
				var errors []ValidationError
				for _, err := range err.(validator.ValidationErrors) {
					errors = append(errors, ValidationError{
						Field:   err.Field(),
						Message: getErrorMsg(err),
					})
				}
				return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails(
					fiber.StatusBadRequest,
					"Bad Request",
					"Query parameter validation failed",
					errors,
				))
			}

			c.Locals("requestQuery", queryInstance)
		}

		if bodySchema != nil {
			bodyInstance := utils.CreateNew(bodySchema)

			if err := c.Bind().Body(bodyInstance); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(response.Error(
					fiber.StatusBadRequest,
					"Bad Request",
					"Invalid request body",
					nil,
				))
			}

			if err := validate.Struct(bodyInstance); err != nil {
				var errors []ValidationError
				for _, err := range err.(validator.ValidationErrors) {
					errors = append(errors, ValidationError{
						Field:   err.Field(),
						Message: getErrorMsg(err),
					})
				}
				return c.Status(fiber.StatusBadRequest).JSON(response.ErrorWithDetails(
					fiber.StatusBadRequest,
					"Bad Request",
					"Request body validation failed",
					errors,
				))
			}

			c.Locals("requestBody", bodyInstance)
		}

		return c.Next()
	}
}

func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Should be greater than " + err.Param()
	case "max":
		return "Should be less than " + err.Param()
	case "oneof":
		return "Should be one of " + err.Param()
	}
	return "Invalid value"
}
