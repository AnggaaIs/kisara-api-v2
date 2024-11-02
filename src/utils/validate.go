package utils

import (
	"kisara/src/config"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(config.AppConfig.JwtKey), nil
	})
}
