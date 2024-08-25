package icloud

import (
	"testing"
	"time"
)

func TestCacheSetAndGet(t *testing.T) {
    cache := NewCache()

    cache.Set("key1", "value1", 5*time.Second)
    value, found := cache.Get("key1")
    if !found || value != "value1" {
        t.Errorf("Expected to find value1, but got %v", value)
    }

    cache.Set("key2", "value2", 1*time.Millisecond)
    time.Sleep(2 * time.Millisecond)
    value, found = cache.Get("key2")
    if found || value != nil {
        t.Errorf("Expected value2 to be expired, but got %v", value)
    }

    cache.Set("key1", "newValue1", 5*time.Second)
    value, found = cache.Get("key1")
    if !found || value != "newValue1" {
        t.Errorf("Expected to find newValue1, but got %v", value)
    }
}
