package controller

import (
	"fmt"
	"kisara/src/models"
	"kisara/src/models/validation"
	"kisara/src/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
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
		requestQuery := c.Locals("requestQuery").(*validation.MessageBodyGet)
		linkID := c.Params("link_id")
		sortBy := requestQuery.SortBy
		page := requestQuery.Page
		limit := requestQuery.Limit

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		offset := (page - 1) * limit

		if sortBy == "" {
			sortBy = "desc"
		}

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
			ID             string    `json:"id"`
			MessageContent string    `json:"message_content"`
			CreatedAt      time.Time `json:"created_at"`
			LikeByCreator  bool      `json:"like_by_creator"`
		}

		var comments []Comment
		var totalComments int64

		if err := db.Model(&models.Comment{}).Where("user_email = ?", user.Email).Count(&totalComments).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to count comments",
				err,
			))
		}

		if err := db.Model(&models.Comment{}).
			Where("user_email = ?", user.Email).
			Order("created_at " + strings.ToUpper(sortBy)).
			Offset(offset).
			Limit(limit).
			Find(&comments).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to retrieve comments",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Comments retrieved successfully",
			map[string]interface{}{
				"page":          page,
				"limit":         limit,
				"total_records": totalComments,
				"total_pages":   (totalComments + int64(limit) - 1) / int64(limit),
				"author": map[string]interface{}{
					"name":        user.Name,
					"profile_url": user.ProfileURL,
				},
				"comments": comments,
			},
		))
	}
}

func HandleDeleteMessage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		if user.Email != email {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"User not authorized to delete message",
				nil,
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

func HandleReplyMessagePost(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")
		content := c.Locals("requestBody").(*validation.MessageBodyPost).MessageContent
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		if user.Email != email {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"User not authorized to reply",
				nil,
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

		reply := models.ReplyComment{
			MessageContent: content,
		}

		if err := db.Model(&comment).Association("ReplyComments").Append(&reply); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to post reply",
				err,
			))
		}

		return c.Status(fiber.StatusCreated).JSON(response.Success(
			fiber.StatusCreated,
			"Created",
			"Reply posted successfully",
			map[string]interface{}{
				"message_id":       messageID,
				"reply_message_id": reply.ID,
				"message_content":  content,
			},
		))
	}
}

func HandleReplyMessageGet(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")
		requestQuery := c.Locals("requestQuery").(*validation.MessageBodyGet)
		page := requestQuery.Page
		limit := requestQuery.Limit
		sortBy := requestQuery.SortBy

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		offset := (page - 1) * limit

		if sortBy == "" {
			sortBy = "desc"
		}

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

		type ReplyComments struct {
			ID             string `json:"id"`
			MessageContent string `json:"message_content"`
			CreatedAt      string `json:"created_at"`
		}

		var replies []ReplyComments
		var totalReplies int64

		if err := db.Model(&models.ReplyComment{}).
			Where("parent_id = ?", messageID).
			Count(&totalReplies).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to count replies",
				err,
			))
		}

		if err := db.Where("parent_id = ?", messageID).
			Order("created_at " + strings.ToUpper(sortBy)).
			Offset(offset).
			Limit(limit).
			Find(&replies).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to retrieve replies",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Replies retrieved successfully",
			map[string]interface{}{
				"page":          page,
				"limit":         limit,
				"total_records": totalReplies,
				"total_pages":   (totalReplies + int64(limit) - 1) / int64(limit),
				"author": map[string]interface{}{
					"name":        user.Name,
					"profile_url": user.ProfileURL,
				},
				"replies": replies,
			},
		))
	}
}

func HandleDeleteReplyMessage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")
		replyID := c.Params("reply_id")
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		if user.Email != email {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"User not authorized to delete reply",
				nil,
			))
		}

		var reply models.ReplyComment
		if err := db.Where("id = ? AND parent_id = ?", replyID, messageID).First(&reply).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("Reply with id %s not found", replyID),
				err,
			))
		}

		if err := db.Delete(&reply).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to delete reply",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Reply deleted successfully",
			nil,
		))
	}
}

func HandleLikeMessage(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		linkID := c.Params("link_id")
		messageID := c.Params("message_id")
		token := c.Locals("token").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		var user models.User
		if err := db.Where("link_id = ?", linkID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("User with link_id %s not found", linkID),
				err,
			))
		}

		if user.Email != email {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(
				fiber.StatusUnauthorized,
				"Unauthorized",
				"User not authorized to like message",
				nil,
			))
		}

		var comment models.Comment
		if err := db.Where("id = ? AND user_email = ?", messageID, user.Email).First(&comment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(
				fiber.StatusNotFound,
				"Not Found",
				fmt.Sprintf("Message with id %s not found", messageID),
				err,
			))
		}

		comment.LikeByCreator = !comment.LikeByCreator

		if err := db.Save(&comment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
				fiber.StatusInternalServerError,
				"Internal Server Error",
				"Failed to update like status",
				err,
			))
		}

		return c.Status(fiber.StatusOK).JSON(response.Success(
			fiber.StatusOK,
			"Success",
			"Like status updated successfully",
			map[string]interface{}{
				"message_id":      comment.ID,
				"like_by_creator": comment.LikeByCreator,
			},
		))
	}
}
