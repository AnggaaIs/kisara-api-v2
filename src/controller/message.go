package controller

import (
	"fmt"
	"kisara/src/models"
	"kisara/src/models/validation"
	"kisara/src/response"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func HandleMessagePost(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		link_id := c.Params("link_id")
		content := c.Locals("requestBody").(*validation.MessageBodyPost)

		var user models.User
		if err := db.Where("link_id = ?", link_id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    fmt.Sprintf("User with link_id %s not found", link_id),
			})
		}

		comment := models.Comment{
			MessageContent: content.MessageContent,
		}

		if err := db.Model(&user).Association("Comments").Append(&comment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to post message",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response.DataResponse{
			GeneralResponse: response.GeneralResponse{StatusCode: fiber.StatusCreated,
				Name:    "Created",
				Message: "Message posted successfully",
			},
			Data: map[string]interface{}{
				"id":              comment.ID,
				"message_content": comment.MessageContent,
			},
		})

	}
}

func HandleMessageGet(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		link_id := c.Params("link_id")

		var user models.User
		if err := db.Where("link_id = ?", link_id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    fmt.Sprintf("User with link_id %s not found", link_id),
			})
		}

		type Comment struct {
			ID             string `json:"id"`
			MessageContent string `json:"message_content"`
			CreatedAt      string `json:"created_at"`
		}

		var comments []Comment

		if err := db.Model(&user).Select("id", "message_content", "created_at").Association("Comments").Find(&comments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to retrieve messages",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.DataResponse{
			GeneralResponse: response.GeneralResponse{
				StatusCode: fiber.StatusOK,
				Name:       "Success",
				Message:    "Messages retrieved successfully",
			},
			Data: comments,
		})
	}
}

func HandleDeleteMessage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {

		linkID := c.Params("link_id")
		messageID := c.Params("message_id")

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    fmt.Sprintf("User with link_id %s not found", linkID),
			})
		}

		var comment models.Comment
		if err := db.Where("id = ?", messageID).First(&comment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusNotFound,
				Name:       "Not Found",
				Message:    fmt.Sprintf("Message with id %s not found", messageID),
			})
		}

		if err := db.Delete(&comment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusInternalServerError,
				Name:       "Internal Server Error",
				Message:    "Failed to delete message",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.GeneralResponse{
			StatusCode: fiber.StatusOK,
			Name:       "Success",
			Message:    "Message deleted successfully",
		})
	}
}
