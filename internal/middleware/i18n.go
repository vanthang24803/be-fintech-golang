package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// I18nMiddleware extracts desired language from headers
func I18nMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Priority 1: Custom header x-lang
		lang := c.Get("x-lang")

		// Priority 2: Standard Accept-Language
		if lang == "" {
			acceptLang := c.Get("Accept-Language")
			if acceptLang != "" {
				// Take first preferred language (e.g. vi-VN,vi;q=0.9 -> vi)
				parts := strings.Split(acceptLang, ",")
				if len(parts) > 0 {
					lang = strings.Split(parts[0], "-")[0]
				}
			}
		}

		// Normalize & Fallback
		lang = strings.ToLower(lang)
		if lang == "" || (lang != "en" && lang != "vi") {
			lang = "en"
		}

		c.Locals("lang", lang)
		return c.Next()
	}
}
