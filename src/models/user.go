package models

import (
	"time"
)

type Role string

const (
	USER    Role = "USER"
	PARTNER Role = "PARTNER"
	BOT     Role = "BOT"
	ADMIN   Role = "ADMIN"
)

type User struct {
	Email      string    `gorm:"primaryKey;unique;not null" json:"email"`
	Name       string    `gorm:"unique;not null" json:"name"`
	LinkID     string    `gorm:"unique;not null" json:"link_id"`
	Role       Role      `gorm:"default:USER" json:"role"`
	ProfileURL *string   `gorm:"null" json:"profile_url"`
	Comments   []Comment `gorm:"foreignKey:UserEmail" json:"comments"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
	HD            string `json:"hd,omitempty"`
	Profile       string `json:"profile,omitempty"`
	Gender        string `json:"gender,omitempty"`
}
