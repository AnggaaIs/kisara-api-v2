package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// server
	ServerPort int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// google oauth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	//key
	JwtKey string
}

var AppConfig *Config

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AppConfig = &Config{
		ServerPort: viper.GetInt("PORT"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),

		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),

		JwtKey: viper.GetString("JWT_KEY"),
	}
}
