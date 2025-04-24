package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/matthewmolinar/tldr/pkg/extract"
	"github.com/matthewmolinar/tldr/pkg/validate"
)

// SummarizeReq represents the request payload for the summarize endpoint
type SummarizeReq struct {
	URL string `json:"url"`
}

// SummarizeResp represents the response from the summarize endpoint
type SummarizeResp struct {
	Headline string   `json:"headline"`
	Bullets  []string `json:"bullets"`
}

// handleSummarize handles article summarization requests
func handleSummarize(c *fiber.Ctx) error {
	// Parse and validate request
	var req SummarizeReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := validate.ValidateURL(req.URL, nil); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid URL format")
	}

	// Extract article text
	text, err := extract.Extract(req.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "failed to extract article content")
	}

	// Generate summary using LLM
	headline, bullets, err := llmClient.Summarize(text)
	if err != nil {
		// Log the actual error
		log.Printf("LLM error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate summary: " + err.Error())
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(SummarizeResp{
		Headline: headline,
		Bullets:  bullets,
	})
}
