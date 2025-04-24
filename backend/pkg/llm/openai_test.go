package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	request  *http.Request
	response *http.Response
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.request = req
	return m.response, nil
}

func TestNewClient(t *testing.T) {
	t.Run("fails without API key", func(t *testing.T) {
		os.Unsetenv("OPENAI_API_KEY")
		_, err := NewClient()
		assert.Error(t, err)
	})

	t.Run("succeeds with API key", func(t *testing.T) {
		os.Setenv("OPENAI_API_KEY", "test-key")
		defer os.Unsetenv("OPENAI_API_KEY")
		
		client, err := NewClient()
		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestClient_Summarize(t *testing.T) {
	// Mock response from OpenAI
	mockResp := openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{
				Message: openai.ChatCompletionMessage{
					Content: "Test Headline",
				},
			},
		},
	}
	respBody, err := json.Marshal(mockResp)
	require.NoError(t, err)

	mock := &mockTransport{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer(respBody)),
		},
	}

	config := openai.DefaultConfig("test-key")
	config.HTTPClient = &http.Client{Transport: mock}
	client := &Client{
		Client: openai.NewClientWithConfig(config),
	}

	t.Run("sends correct prompt format", func(t *testing.T) {
		headline, bullets, err := client.Summarize("Test article content")
		require.NoError(t, err)
		assert.Equal(t, "Test Headline", headline)
		assert.Empty(t, bullets, "bullets should be empty for now")

		// Verify request was properly formed
		assert.NotNil(t, mock.request)
		
		// Read and verify request body
		var reqBody openai.ChatCompletionRequest
		err = json.NewDecoder(mock.request.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify system prompt
		require.GreaterOrEqual(t, len(reqBody.Messages), 2)
		assert.Equal(t, openai.ChatMessageRoleSystem, reqBody.Messages[0].Role)
		assert.Contains(t, reqBody.Messages[0].Content, "master headline writer")
		
		// Verify user content
		assert.Equal(t, openai.ChatMessageRoleUser, reqBody.Messages[1].Role)
		assert.Equal(t, "Test article content", reqBody.Messages[1].Content)

		// Verify temperature
		assert.Equal(t, float32(0.5), reqBody.Temperature)
	})
}
