package validate

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	responseSize int64
	failRequest  bool
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.failRequest {
		return nil, fmt.Errorf("failed to fetch")
	}
	return &http.Response{
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(strings.NewReader("")),
		ContentLength: m.responseSize,
	}, nil
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		responseSize int64
		failRequest  bool
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid https url with small content",
			url:          "https://example.com",
			responseSize: 1000,
			wantErr:      false,
		},
		{
			name:        "invalid scheme (http)",
			url:         "http://example.com",
			wantErr:     true,
			errContains: "must use HTTPS",
		},
		{
			name:        "invalid url format",
			url:         "not-a-url",
			wantErr:     true,
			errContains: "invalid URL format",
		},
		{
			name:         "content too large",
			url:          "https://example.com",
			responseSize: maxContentLengthBytes + 1,
			wantErr:      true,
			errContains:  "content too large",
		},
		{
			name:         "unreachable url",
			url:          "https://this-should-not-exist.test",
			failRequest:  true,
			wantErr:      true,
			errContains:  "failed to fetch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{
				Transport: &mockTransport{
					responseSize: tt.responseSize,
					failRequest:  tt.failRequest,
				},
			}

			err := ValidateURL(tt.url, client)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), tt.errContains),
					"error should contain %q, got %q", tt.errContains, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
