# Changelog

## [0.1.0] - 2025-04-24

### Added
- Initial Go project setup with proper module structure
- Basic Fiber server implementation with health check endpoint
- Configured server timeouts (5s) and environment-based port
- Unit tests for health check endpoint using httptest
- Added core dependencies:
  - github.com/gofiber/fiber/v2
  - github.com/stretchr/testify for testing

### [B2] - 2025-04-24
- Added `/api/summarize` endpoint with dummy response
- Implemented request logging middleware
- Added integration tests for server startup and routes
- Established test patterns with setupTestApp helper

### [B3] - 2025-04-24
- Added URL validation package with strict requirements:



### [B4] - 2025-04-24
- Fixed content extraction using go-shiori/go-readability
- Added proper error handling for article parsing
- Implemented 8KB content size limit
- Added integration tests with sample article
  - HTTPS scheme enforcement
  - 20KB size limit for HTML content
- Integrated validation into summarize endpoint
- Added comprehensive table-driven tests with HTTP mocking


### [B5] - 2025-04-24
- Implemented OpenAI client wrapper in pkg/llm
- Added gpt-3.5-turbo integration with system prompt for headline writing
- Configured temperature=0.5 as per PRD requirements
- Added unit tests with mock HTTP client for OpenAI API calls
- Added environment-based API key configuration
