package clients

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGarminClient(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		rateLimitAllow bool
		httpStatusCode int
		expectedError  string
	}{
		{
			name:           "Message length exceeds limit",
			message:        "This is a very long message that exceeds the limit  a very long message that exceeds the limit This is a very long message that exceeds the limit  a very long message that exceeds the limit",
			rateLimitAllow: true,
			httpStatusCode: http.StatusOK,
			expectedError:  "message length exceeds limit",
		},
		{
			name:           "Rate limit exceeded",
			message:        "Hello, World!",
			rateLimitAllow: false,
			httpStatusCode: http.StatusOK,
			expectedError:  "rate limit exceeded",
		},
		{
			name:           "Non-OK HTTP status code",
			message:        "Hello, World!",
			rateLimitAllow: true,
			httpStatusCode: http.StatusInternalServerError,
			expectedError:  "failed to send message, status code: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the real rate limiter
			rateLimiter := NewRateLimiter(1, time.Minute)

			// Mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.httpStatusCode)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			}))
			defer server.Close()

			client := &GarminClient{
				httpClient:  server.Client(),
				address:     server.URL,
				imei:        "test-imei",
				rateLimiter: rateLimiter,
			}

			// Simulate rate limiting
			if !tt.rateLimitAllow {
				rateLimiter.Allow("test-imei")
				rateLimiter.Allow("test-imei")
			}

			// Call the method under test
			err := client.SendMessage(tt.name, tt.message)

			// Check the expected error
			if err == nil || err.Error() != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}