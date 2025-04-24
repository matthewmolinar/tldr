package main

import (
	"github.com/gofiber/fiber/v2"
)

// SummarizeReq represents the request payload for the summarize endpoint
type SummarizeReq struct {
	URL string `json:"url"`
}

// summarizeHandler handles article summarization requests
func summarizeHandler(c *fiber.Ctx) error {
	// Return dummy JSON for now
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"headline": "stub",
		"bullets":  []string{},
	})
}
