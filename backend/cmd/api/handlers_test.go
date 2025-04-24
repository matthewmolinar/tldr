package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	
	api := app.Group("/api")
	api.Post("/summarize", summarizeHandler)
	return app
}

func TestServerStartup(t *testing.T) {
	app := setupTestApp()
	assert.NotNil(t, app, "Server should initialize")

	// Test route registration
	req := httptest.NewRequest("GET", "/non-existent", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Unknown route should return 404")
}

func TestSummarizeHandler(t *testing.T) {
	app := setupTestApp()

	t.Run("returns 201 with dummy JSON", func(t *testing.T) {
		reqBody := `{"url":"https://example.com"}`
		req := httptest.NewRequest("POST", "/api/summarize", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"headline":"stub","bullets":[]}`, string(body))
	})
}
