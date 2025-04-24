package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestErrorMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorMiddleware())

	// Create test handler that returns an error
	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.ErrUnprocessableEntity
	})

	// Make request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check status code
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	// Check response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Verify error message format
	assert.Contains(t, result, "error")
	assert.Equal(t, fiber.ErrUnprocessableEntity.Error(), result["error"])
}
