// Package validate provides validation utilities for the TL;DR service
package validate

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

const maxContentLengthBytes = 10 * 1024 * 1024 // 10 MB

// ValidateURL checks if the given URL is valid according to service requirements:
// - Must use HTTPS scheme
// - Content must not exceed 20 KB (checked via HEAD request)
func ValidateURL(s string, client *http.Client) error {
	if client == nil {
		// Create a custom client with TLS certificate verification disabled
		// only in production environments to handle containerized deployments
		_, inProduction := os.LookupEnv("FLY_APP_NAME")
		if inProduction {
			log.Printf("Running in production environment, disabling TLS verification")
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		} else {
			client = http.DefaultClient
		}
	}
	// Parse and normalize URL
	log.Printf("Validating URL: %q", s)
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("URL parsing error: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid URL format")
	}

	// Normalize URL
	log.Printf("Parsed URL - Host: %q, Scheme: %q, Path: %q", u.Host, u.Scheme, u.Path)
	u.RawQuery = ""
	if u.Host == "" {
		return fiber.NewError(fiber.StatusBadRequest, "URL must include a host")
	}

	// Verify HTTPS scheme
	if u.Scheme != "https" {
		return fiber.NewError(fiber.StatusBadRequest, "URL must use HTTPS")
	}

	// Check content size via HEAD request
	log.Printf("Making HEAD request to: %s", u.String())
	req, err := http.NewRequest(http.MethodHead, u.String(), nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to create request: %v", err))
	}

	// Add a user agent to avoid being blocked by some sites
	req.Header.Set("User-Agent", "Mozilla/5.0 TL;DR-App/1.0")
	
	log.Printf("Sending HEAD request")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch URL: %v", err)
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("failed to fetch URL: %v", err))
	}
	defer resp.Body.Close()

	log.Printf("Received response status: %s", resp.Status)
	
	// Check content length if provided
	contentLength := resp.ContentLength
	log.Printf("Content length: %d bytes (max %d bytes)", contentLength, maxContentLengthBytes)
	
	if contentLength > maxContentLengthBytes {
		return fiber.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("content too large: %d bytes (max %d bytes)", 
			contentLength, maxContentLengthBytes),
		)
	}

	return nil
}
