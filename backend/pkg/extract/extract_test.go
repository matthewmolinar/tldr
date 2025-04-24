package extract

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract(t *testing.T) {
	// Load test HTML file
	htmlData, err := os.ReadFile("testdata/article.html")
	require.NoError(t, err)

	// Create test server that serves our test HTML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlData)
	}))
	defer ts.Close()

	// Test extraction
	content, err := Extract(ts.URL)
	require.NoError(t, err)

	// Assert content is within size limit
	assert.LessOrEqual(t, len(content), maxBytes)

	// Assert content contains expected text (update this based on your test article)
	assert.Contains(t, content, "This is the main content")
}

func TestExtract_Errors(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "invalid URL",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "non-existent URL",
			url:     "http://localhost:1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Extract(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
