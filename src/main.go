package main

import (
	"context"
	"fmt"
	"kisara/src/config"
	"kisara/src/database"
	"kisara/src/response"
	"kisara/src/router"
	"kisara/src/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := setupFiber()
	db := setupDatabase()
	defer closeDatabase(db)
	setupRouter(app, db)

	go startFiberServer(app)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		utils.Log.Info("Shutdown signal received")
		cancel()
	}()

	// Panggil fungsi handleGracefulShutdown
	handleGracefulShutdown(ctx, app)
}

func setupFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		Format:   "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeZone: "Asia/Jakarta",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        100,
		Expiration: 60 * time.Second,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(response.GeneralResponse{
				StatusCode: fiber.StatusTooManyRequests,
				Name:       "Too Many Requests",
				Message:    "You’ve reached the request limit. Kick back for a moment and come back later.",
			})
		},
	}))

	return app
}

func startFiberServer(app *fiber.App) {
	utils.Log.Info("Starting server...")

	PORT := config.AppConfig.ServerPort
	serverPortS := fmt.Sprintf(":%d", PORT)

	if err := app.Listen(serverPortS); err != nil {
		utils.Log.Fatal("Error when starting server", err)
	}

}

func setupDatabase() *gorm.DB {
	utils.Log.Info("Setting up database...")

	db := database.Connect()

	return db
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()

	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func setupRouter(app *fiber.App, db *gorm.DB) {
	utils.Log.Info("Setting up router...")

	router.SetupRoutes(app, db)
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App) {
	<-ctx.Done()

	utils.Log.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		utils.Log.Error("Server shutdown failed: ", err)
	}

	utils.Log.Info("Server stopped gracefully")
}
