package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID             uuid.UUID      `gorm:"primaryKey;unique;default:uuid_generate_v4()" json:"id"`
	MessageContent string         `json:"message_content"`
	UserEmail      string         `gorm:"index;not null" json:"user_email"`
	ReplyComments  []ReplyComment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"reply_comments"`
	LikeByCreator  bool           `gorm:"default:false" json:"like_by_creator"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}
