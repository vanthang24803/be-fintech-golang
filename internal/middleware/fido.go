package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtUtil "github.com/maynguyen24/sever/pkg/jwt"
	"github.com/maynguyen24/sever/pkg/response"
)

// FIDOMiddleware ensures that the current session has been verified via Biometrics/FIDO.
// This should be applied AFTER AuthRequired middleware.
func FIDOMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract claims from context (populated by AuthRequired)
		userClaims, ok := c.Locals("user").(*jwtUtil.TokenClaims)
		if !ok || userClaims == nil {
			return response.Error(c, fiber.StatusUnauthorized, 4010, "common.unauthorized", "Missing user claims")
		}

		// Check if FIDO is verified
		if !userClaims.FIDOVerified {
			// 403 Forbidden with a specific code to trigger FIDO prompt on Frontend
			return response.Error(c, fiber.StatusForbidden, 4031, "security.fido_required", "Step-up authentication needed")
		}

		return c.Next()
	}
}
