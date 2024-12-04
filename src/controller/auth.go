package controller

import (
	"errors"
	"kisara/src/config"
	"kisara/src/models"
	"kisara/src/models/validation"
	"kisara/src/response"
	"kisara/src/utils"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func HandleGoogleURL(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Generate random state
		state, err := utils.GenerateRandomState(21)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to generate random state",
					err,
				),
			)
		}

		// Return success response with generated URL
		return c.Status(fiber.StatusOK).JSON(
			response.Success(
				fiber.StatusOK,
				"Success",
				"Google URL generated successfully",
				map[string]interface{}{
					"url": config.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce),
				},
			),
		)
	}
}

func HandleGoogleCallback(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract request body
		query := c.Locals("requestBody").(*validation.AuthGoogleCallbackBody)

		// Exchange code with token
		token, err := config.GoogleOAuthConfig.Exchange(c.Context(), query.Code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to exchange code with token",
					err,
				),
			)
		}

		// Get user info from Google
		getUserInfo, err := config.GoogleOAuthConfig.Client(c.Context(), token).Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to get user info",
					err,
				),
			)
		}
		defer getUserInfo.Body.Close()

		// Decode user info response
		var userInfo models.GoogleUserInfo
		if err := json.NewDecoder(getUserInfo.Body).Decode(&userInfo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to decode user info",
					err,
				),
			)
		}

		// Generate a unique link ID for the user
		linkID, err := utils.GenerateUniqueLinkID(db, 5)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to generate unique link ID",
					err,
				),
			)
		}

		// Create JWT token with user claims
		claims := jwt.MapClaims{
			"email":      userInfo.Email,
			"sub":        userInfo.Sub,
			"picture":    userInfo.Picture,
			"name":       userInfo.Name,
			"link_id":    linkID,
			"time_login": time.Now(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessToken, err := jwtToken.SignedString([]byte(config.AppConfig.JwtKey))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to generate access token",
					err,
				),
			)
		}

		// Check if the user exists in the database
		errDb := db.First(&models.User{}, "email = ?", userInfo.Email).Error
		if errors.Is(errDb, gorm.ErrRecordNotFound) {
			// Create a new user if not found
			if err := db.Create(&models.User{
				Email:      userInfo.Email,
				Name:       userInfo.Name,
				LinkID:     linkID,
				ProfileURL: &userInfo.Picture,
			}).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					response.Error(
						fiber.StatusInternalServerError,
						"Internal Server Error",
						"Failed to create user",
						err,
					),
				)
			}
		} else if errDb == nil {
			// Update existing user information
			if err := db.Model(&models.User{}).Updates(models.User{
				Name:       userInfo.Name,
				ProfileURL: &userInfo.Picture,
			}).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					response.Error(
						fiber.StatusInternalServerError,
						"Internal Server Error",
						"Failed to update user",
						err,
					),
				)
			}
		} else {
			// Handle unexpected database error
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.Error(
					fiber.StatusInternalServerError,
					"Internal Server Error",
					"Failed to process user data",
					errDb,
				),
			)
		}

		// Success response with the generated access token
		return c.Status(fiber.StatusOK).JSON(
			response.Success(
				fiber.StatusOK,
				"Success",
				"Google callback success",
				map[string]interface{}{
					"access_token": accessToken,
				},
			),
		)
	}
}
