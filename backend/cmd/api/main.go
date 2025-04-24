package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/matthewmolinar/tldr/pkg/llm"
)

// Global LLM client for reuse
var llmClient llm.Summarizer

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize LLM client
	// Debug: Print first 10 chars of API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) > 10 {
		log.Printf("Using API key starting with: %s...", apiKey[:10])
	} else {
		log.Printf("API key is too short or empty")
	}

	var err error
	llmClient, err = llm.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize LLM client: %v", err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Add logger middleware
	app.Use(logger.New())

	// Health check endpoint
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// API routes
	api := app.Group("/api")
	api.Post("/summarize", handleSummarize)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
