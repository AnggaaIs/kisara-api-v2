package utils

import (
	crand "crypto/rand"
	"encoding/base64"
	"kisara/src/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func GenerateRandomState(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := crand.Read(bytes)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	return state[:length], nil
}

func GenerateLinkID(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	id := make([]byte, length)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

func GenerateUniqueLinkID(db *gorm.DB, length int) (string, error) {
	for {
		linkID := GenerateLinkID(length)
		var count int64
		if err := db.Model(&models.User{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return linkID, nil
		}
	}
}
