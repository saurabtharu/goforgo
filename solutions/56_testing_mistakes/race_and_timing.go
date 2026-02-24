// race_and_timing.go
// Solution for Mistakes #83, #86, #87
//
// Fixed:
// 1. Counter uses sync.Mutex for safe concurrent access (Mistake #83)
// 2. Cache accepts a Clock interface for testable time-dependent code (Mistake #87)

package main

import (
	"fmt"
	"sync"
	"time"
)

// Counter tracks a count that can be incremented concurrently.
// FIXED: Uses sync.Mutex to prevent data races.
type Counter struct {
	mu    sync.Mutex
	count int
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// Clock is an interface for getting the current time.
// This allows injecting a fake clock in tests (Mistake #87).
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the actual system time.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

// CacheEntry holds a value with an expiration time.
type CacheEntry struct {
	Value     string
	ExpiresAt time.Time
}

// Cache stores string values with time-based expiration.
// FIXED: Uses Clock interface instead of time.Now() directly.
type Cache struct {
	entries map[string]CacheEntry
	ttl     time.Duration
	clock   Clock
}

// NewCache creates a cache with the given TTL and clock.
// FIXED: Accepts a Clock parameter for testability.
func NewCache(ttl time.Duration, clock Clock) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
		clock:   clock,
	}
}

func (c *Cache) Set(key, value string) {
	// FIXED: Uses injected clock instead of time.Now()
	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: c.clock.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}
	// FIXED: Uses injected clock instead of time.Now()
	if c.clock.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		return "", false
	}
	return entry.Value, true
}

func main() {
	fmt.Println("Run 'go test -race -v' to execute the tests with race detection.")
}
