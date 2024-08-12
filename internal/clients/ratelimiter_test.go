package clients

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limit := 3
	interval := 1 * time.Second
	rl := NewRateLimiter(limit, interval)

	key := "test-key"

	for i := 0; i < limit; i++ {
		if !rl.Allow(key) {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	if rl.Allow(key) {
		t.Errorf("expected request to be blocked after reaching the limit")
	}

	time.Sleep(interval)

	if !rl.Allow(key) {
		t.Errorf("expected request to be allowed after interval has passed")
	}
}
