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
		linkID := c.Params("link_id")
		content := c.Locals("requestBody").(*validation.MessageBodyPost)

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		comment := models.Comment{
			MessageContent: content.MessageContent,
		}

		if err := db.Model(&user).Association("Comments").Append(&comment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to post message",
				err,
			))
		}

		return c.Status(fiber.StatusCreated).JSON(response.Success(
			fiber.StatusCreated,
			"Created",
			"Message posted successfully",
			map[string]interface{}{
				"id":              comment.ID,
				"message_content": comment.MessageContent,
			},
		))
	}
}

func HandleMessageGet(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		type Comment struct {
			ID             string `json:"id"`
			MessageContent string `json:"message_content"`
			CreatedAt      string `json:"created_at"`
		}

		var comments []Comment
		if err := db.Model(&user).Select("id", "message_content", "created_at").Association("Comments").Find(&comments); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to retrieve messages",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Messages retrieved successfully",
			comments,
		))
	}
}

func HandleDeleteMessage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		var comment models.Comment
		if err := db.Where("id = ?", messageID).First(&comment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("Message with id %s not found", messageID),
				err,
			))
		}

		if err := db.Delete(&comment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to delete message",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Message deleted successfully",
			nil,
		))
	}
}
