package garmin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGarminClient_MessageLengthExceedsLimit(t *testing.T) {
	rateLimiter := NewRateLimiter(1, time.Minute)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient:  server.Client(),
		address:     server.URL,
		imei:        "test-imei",
		rateLimiter: rateLimiter,
	}

	message := "This is a very long message that exceeds the limit a very long message that exceeds the limit This is a very long message that exceeds the limit a very long message that exceeds the limit"

	err := client.SendMessage("test-sender", message)
	expectedError := "message length exceeds limit"
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestGarminClient_RateLimitExceeded(t *testing.T) {
	rateLimiter := NewRateLimiter(1, time.Minute)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient:  server.Client(),
		address:     server.URL,
		imei:        "test-imei",
		rateLimiter: rateLimiter,
	}

	// Exceed the rate limit
	rateLimiter.Allow("test-imei")
	rateLimiter.Allow("test-imei")

	err := client.SendMessage("test-sender", "Hello, World!")
	expectedError := "rate limit exceeded"
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestGarminClient_NonOKHTTPStatusCode(t *testing.T) {
	rateLimiter := NewRateLimiter(1, time.Minute)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		httpClient:  server.Client(),
		address:     server.URL,
		imei:        "test-imei",
		rateLimiter: rateLimiter,
	}

	err := client.SendMessage("test-sender", "Hello, World!")
	expectedError := "failed to send message, status code: 500"
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
