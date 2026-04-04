package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/logger"
	"go.uber.org/zap"
)

// RequestLogger intercepts HTTP requests and logs their details using Zap
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		logger.Log.Info("Request handled",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
		)

		return err
	}
}
