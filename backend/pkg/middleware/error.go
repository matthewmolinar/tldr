package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorMiddleware converts errors to JSON responses with appropriate status codes
func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue stack
		err := c.Next()
		if err == nil {
			return nil
		}

		// Get fiber's error status code if it exists
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Return JSON error response
		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
