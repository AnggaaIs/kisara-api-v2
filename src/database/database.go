package database

import (
	"database/sql"
	"fmt"
	"kisara/src/config"
	"kisara/src/models"
	"kisara/src/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	sslmode := "disable"
	if config.AppConfig.Env == "production" {
		sslmode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		config.AppConfig.DBHost, config.AppConfig.DBUser, config.AppConfig.DBPassword,
		config.AppConfig.DBName, config.AppConfig.DBPort, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		utils.Log.Fatal("Error when connecting to database", err)
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Fatal("Error when setting up database", errDB)
	}

	// Set optimal connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	if err := pingDatabase(sqlDB); err != nil {
		utils.Log.Fatal("Koneksi ke database gagal:", err)
	}

	utils.Log.Info("Database connection established successfully")

	// Enable UUID Extension
	enableUUIDExtension(db)

	// Migrate models to database
	if err := db.AutoMigrate(&models.User{}, &models.Comment{}, &models.ReplyComment{}); err != nil {
		utils.Log.Fatal("Error when migrating database", err)
	}

	utils.Log.Info("Database migration successful")

	return db
}

func pingDatabase(sqlDB *sql.DB) error {
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed ping to database: %w", err)
	}
	return nil
}

func enableUUIDExtension(db *gorm.DB) {
	// Checking if the UUID extension already exists
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		utils.Log.Fatalf("failed to enable uuid-ossp extension: %v", err)
	}
}
