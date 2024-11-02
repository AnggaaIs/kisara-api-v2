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
	"gorm.io/gorm"
)

func HandleGoogleURL(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		state, err := utils.GenerateRandomState(21)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to generate random state",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.DataResponse{
			GeneralResponse: response.GeneralResponse{
				StatusCode: fiber.StatusOK,
				Name:       "Success",
				Message:    "Google URL generated successfully",
			},
			Data: map[string]interface{}{
				"url": config.GoogleOAuthConfig.AuthCodeURL(state),
			},
		})
	}
}

func HandleGoogleCallback(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		query := c.Locals("requestBody").(*validation.AuthGoogleCallbackBody)

		token, err := config.GoogleOAuthConfig.Exchange(c.Context(), query.Code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to exchange code with token",
			})
		}

		getUserInfo, err := config.GoogleOAuthConfig.Client(c.Context(), token).Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to get user info",
			})
		}
		defer getUserInfo.Body.Close()

		var userInfo models.GoogleUserInfo
		if err := json.NewDecoder(getUserInfo.Body).Decode(&userInfo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to decode user info",
			})
		}

		linkID := utils.GenerateLinkID(7)

		claims := jwt.MapClaims{
			"email":      userInfo.Email,
			"sub":        userInfo.Sub,
			"picture":    userInfo.Picture,
			"name":       userInfo.Name,
			"link_id":    linkID,
			"time_login": time.Now(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessToken, err := jwtToken.SignedString(config.AppConfig.JwtKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to generate access token",
			})
		}

		errDb := db.First(&models.User{}, "email = ?", userInfo.Email).Error

		if errors.Is(errDb, gorm.ErrRecordNotFound) {
			if err := db.Create(&models.User{
				Email:      userInfo.Email,
				Name:       userInfo.Name,
				LinkID:     linkID,
				ProfileURL: &userInfo.Picture,
			}).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
					StatusCode: fiber.StatusInternalServerError,
					Name:       "Internal Server Error",
					Message:    "Failed to create user",
				})
			}
		} else if errDb == nil {
			if err := db.Model(&models.User{}).Updates(models.User{
				Name:       userInfo.Name,
				ProfileURL: &userInfo.Picture,
			}).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
					StatusCode: fiber.StatusInternalServerError,
					Name:       "Internal Server Error",
					Message:    "Failed to update user",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to process user data",
			})
		}

		//success
		return c.Status(fiber.StatusOK).JSON(response.DataResponse{
			GeneralResponse: response.GeneralResponse{
				StatusCode: fiber.StatusOK,
				Name:       "Success",
				Message:    "Google callback success",
			},
			Data: map[string]interface{}{
				"access_token": accessToken,
			},
		})
	}
}
