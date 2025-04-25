package llm

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Summarizer defines the interface for text summarization
type Summarizer interface {
	Summarize(text string) (headline string, bullets []string, err error)
}

// Client wraps the OpenAI client to provide summarization capabilities
type Client struct {
	*openai.Client
}

// NewClient creates a new OpenAI client using the API key from environment
func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}
	
	// Check if running in production environment
	_, inProduction := os.LookupEnv("FLY_APP_NAME")
	
	var client *openai.Client
	if inProduction {
		log.Printf("Running in production environment, disabling TLS verification for OpenAI client")
		// Create a custom HTTP client with TLS certificate verification disabled
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		
		client = openai.NewClientWithConfig(openai.ClientConfig{
			AuthToken: apiKey,
			HTTPClient: httpClient,
		})
	} else {
		client = openai.NewClient(apiKey)
	}
	
	return &Client{client}, nil
}

// Summarize takes an article text and returns a headline and bullet points
func (c *Client) Summarize(text string) (headline string, bullets []string, err error) {
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "You are a master headline writer. Provide a one-sentence headline and 3 bullet " +
						"takeaway points, total < 280 chars.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
			Temperature: 0.5,
		},
	)
	if err != nil {
		return "", nil, err
	}

	// Parse the response into headline and bullets
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", nil, errors.New("no summary generated")
	}

	// Parse response - format is expected to be:
	// Headline: ...
	// - Point 1
	// - Point 2
	// - Point 3
	lines := strings.Split(resp.Choices[0].Message.Content, "\n")
	if len(lines) < 4 {
		return "", nil, errors.New("invalid response format")
	}

	// Extract headline (remove "Headline: " prefix if present)
	headline = strings.TrimPrefix(lines[0], "Headline: ")

	// Extract bullet points
	bullets = make([]string, 0, 3)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			bullets = append(bullets, strings.TrimPrefix(line, "- "))
		}
	}

	return headline, bullets, nil
}
