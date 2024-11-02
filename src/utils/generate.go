package utils

import (
	crand "crypto/rand"
	"encoding/base64"
	"math/rand"
	"time"
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
