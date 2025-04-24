package extract

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-shiori/go-readability"
)

const maxBytes = 8192 // 8 KB limit as per PRD

// Extract fetches the given URL and returns its main content using readability,
// trimmed to a maximum of 8KB.
func Extract(url string) (string, error) {
	// Fetch the page
	log.Printf("Fetching URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", url, err)
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Log response status
	log.Printf("Response status: %s", resp.Status)

	// Parse with readability
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("Read %d bytes from response body", len(body))

	// Parse with readability
	parser := readability.NewParser()
	doc, err := parser.Parse(bytes.NewReader(body), resp.Request.URL)
	if err != nil {
		log.Printf("Failed to parse content: %v", err)
		return "", fmt.Errorf("failed to parse content: %w", err)
	}
	log.Printf("Successfully parsed content with readability")

	// Get content and validate
	content := doc.TextContent
	if content == "" {
		log.Printf("No content extracted from URL %s", url)
		return "", fmt.Errorf("no content extracted from URL")
	}
	log.Printf("Extracted %d characters of content", len(content))

	// Trim if needed
	if len(content) > maxBytes {
		content = content[:maxBytes]
	}

	return content, nil
}
