package models

import (
	"time"

	"github.com/google/uuid"
)

type ReplyComment struct {
	ID             uuid.UUID `gorm:"primaryKey;unique;default:uuid_generate_v4()" json:"id"`
	ParentID       uuid.UUID `gorm:"not null" json:"parent_id"` // Foreign key
	MessageContent string    `json:"message_content"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
