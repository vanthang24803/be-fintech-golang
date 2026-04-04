package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/i18n"
)

// Response struct defines the standard format for all API responses
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Success is a helper to return a unified success response with i18n translation
func Success(c *fiber.Ctx, code int, messageKey string, data any) error {
	lang, _ := c.Locals("lang").(string)
	if lang == "" {
		lang = "en"
	}

	translatedMessage := i18n.Translate(lang, messageKey)

	return c.Status(fiber.StatusOK).JSON(&Response{
		Code:    code,
		Message: translatedMessage,
		Data:    data,
	})
}

// Error is a helper to return a unified error response with i18n translation
func Error(c *fiber.Ctx, httpStatus int, code int, messageKey string, errDetail string) error {
	lang, _ := c.Locals("lang").(string)
	if lang == "" {
		lang = "en"
	}

	translatedMessage := i18n.Translate(lang, messageKey)

	return c.Status(httpStatus).JSON(&Response{
		Code:    code,
		Message: translatedMessage,
		Error:   errDetail,
	})
}
