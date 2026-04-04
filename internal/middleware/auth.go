package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
)

// AuthRequired is an interceptor that validates the JWT Access Token
func AuthRequired(cfg *configs.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization format")
		}

		tokenString := parts[1]

		// Parse and validate the token using our custom claims structure
		token, err := jwt.ParseWithClaims(tokenString, &jwtUtil.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Extract Context
		claims, ok := token.Claims.(*jwtUtil.TokenClaims)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Failed to parse claims")
		}

		// Inject User ID and full claims into Fiber Context locally for downstream handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("user", claims)

		return c.Next()
	}
}
