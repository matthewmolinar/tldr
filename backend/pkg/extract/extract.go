package extract

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-shiori/go-readability"
)

const maxBytes = 8192 // 8 KB limit as per PRD

// Extract fetches the given URL and returns its main content using readability,
// trimmed to a maximum of 8KB.
func Extract(url string) (string, error) {
	// Fetch the page
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Parse with readability
	// Read response body
body, err := io.ReadAll(resp.Body)
if err != nil {
	return "", fmt.Errorf("failed to read response body: %w", err)
}

// Parse with readability
parser := readability.NewParser()
	doc, err := parser.Parse(bytes.NewReader(body), resp.Request.URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse content: %w", err)
	}

	// Get content and trim if needed
	content := doc.TextContent
	if len(content) > maxBytes {
		content = content[:maxBytes]
	}

	return content, nil
}
