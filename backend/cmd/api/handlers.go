package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matthewmolinar/tldr/pkg/validate"
)

// SummarizeReq represents the request payload for the summarize endpoint
type SummarizeReq struct {
	URL string `json:"url"`
}

// summarizeHandler handles article summarization requests
func summarizeHandler(c *fiber.Ctx) error {
	var req SummarizeReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validate.ValidateURL(req.URL, nil); err != nil {
		return err // ValidateURL already returns fiber.Error
	}

	// Return dummy JSON for now
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"headline": "stub",
		"bullets":  []string{},
	})
}
