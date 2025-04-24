package llm

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

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
	return &Client{openai.NewClient(apiKey)}, nil
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
	// TODO: Add proper parsing logic based on the actual response format
	if len(resp.Choices) == 0 {
		return "", nil, errors.New("no completion choices returned")
	}

	// For now, returning the raw response - this will be refined in the next iteration
	content := resp.Choices[0].Message.Content
	return content, []string{}, nil
}
