package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/matthewmolinar/tldr/pkg/llm"
	"github.com/stretchr/testify/assert"
)

// mockLLMClient is a test double that returns canned responses
type mockLLMClient struct{}

func (m *mockLLMClient) Summarize(text string) (string, []string, error) {
	return "Test Headline", []string{"Point 1", "Point 2", "Point 3"}, nil
}

func setupTestApp(client llm.Summarizer) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	
	// Override global client for testing
	llmClient = client
	
	api := app.Group("/api")
	api.Post("/summarize", handleSummarize)
	return app
}

func TestServerStartup(t *testing.T) {
	app := setupTestApp(&mockLLMClient{})
	assert.NotNil(t, app, "Server should initialize")

	// Test route registration
	req := httptest.NewRequest("GET", "/non-existent", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Unknown route should return 404")
}

func TestSummarizeHandler(t *testing.T) {
	app := setupTestApp(&mockLLMClient{})

	t.Run("returns 201 with summary for valid URL", func(t *testing.T) {
		reqBody := `{"url":"https://example.com"}`
		req := httptest.NewRequest("POST", "/api/summarize", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		expected := `{"headline":"Test Headline","bullets":["Point 1","Point 2","Point 3"]}`
		assert.JSONEq(t, expected, string(body))
	})

	t.Run("returns 400 for invalid request body", func(t *testing.T) {
		reqBody := `{invalid json}`
		req := httptest.NewRequest("POST", "/api/summarize", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 422 for invalid URL", func(t *testing.T) {
		reqBody := `{"url":"not-a-url"}`
		req := httptest.NewRequest("POST", "/api/summarize", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})
}
