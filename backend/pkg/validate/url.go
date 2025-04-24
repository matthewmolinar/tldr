// Package validate provides validation utilities for the TL;DR service
package validate

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

const maxContentLengthBytes = 20 * 1024 // 20 KB

// ValidateURL checks if the given URL is valid according to service requirements:
// - Must use HTTPS scheme
// - Content must not exceed 20 KB (checked via HEAD request)
func ValidateURL(s string, client *http.Client) error {
	if client == nil {
		client = http.DefaultClient
	}
	// Parse URL
	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid URL format")
	}

	// Verify HTTPS scheme
	if u.Scheme != "https" {
		return fiber.NewError(fiber.StatusBadRequest, "URL must use HTTPS")
	}

	// Check content size via HEAD request
	req, err := http.NewRequest(http.MethodHead, s, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to create request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to fetch URL")
	}
	defer resp.Body.Close()

	// Check content length if provided
	contentLength := resp.ContentLength
	if contentLength > maxContentLengthBytes {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("content too large: %d bytes (max %d bytes)", 
				contentLength, maxContentLengthBytes),
		)
	}

	return nil
}
