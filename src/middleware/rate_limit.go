package middleware

import (
	"kisara/src/response"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type RateLimitConfig struct {
	Max      int
	Duration int
}

func RateLimitMiddleware(config RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.Max,
		Expiration: time.Duration(config.Duration) * time.Second,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(response.Error(
				fiber.StatusTooManyRequests,
				"Too Many Requests",
				"You’ve reached the request limit. Kick back for a moment and come back later.",
				nil,
			))
		},
	})
}
