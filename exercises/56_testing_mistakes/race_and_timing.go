// race_and_timing.go
// Mistakes #83, #86, #87: Race detection, sleep-based tests, and time abstraction
//
// This file has two problems to fix:
// 1. Counter has a data race (no synchronization) - Mistake #83
// 2. Cache uses time.Now() directly, making it hard to test - Mistake #87
//
// Fix both issues in this file, then fix the tests in race_and_timing_test.go.

package main

import (
	"fmt"
	"time"
)

// Counter tracks a count that can be incremented concurrently.
// BUG: This has a data race - multiple goroutines read/write count
// without synchronization. Fix it by adding a sync.Mutex.
type Counter struct {
	count int
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Increment() {
	// BUG: No synchronization - data race when called from multiple goroutines
	c.count++
}

func (c *Counter) Value() int {
	// BUG: No synchronization - data race when called while incrementing
	return c.count
}

// Clock is an interface for getting the current time.
// This allows injecting a fake clock in tests (Mistake #87).
// TODO: The Cache below uses time.Now() directly. Refactor it to use
// this Clock interface instead, so tests can control time.
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
// BUG: Uses time.Now() directly, making it impossible to test
// expiration without time.Sleep. Refactor to use the Clock interface.
type Cache struct {
	entries map[string]CacheEntry
	ttl     time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
}

func (c *Cache) Set(key, value string) {
	// BUG: Uses time.Now() directly - should use an injected Clock
	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}
	// BUG: Uses time.Now() directly - should use an injected Clock
	if time.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		return "", false
	}
	return entry.Value, true
}

func main() {
	fmt.Println("Run 'go test -race -v' to execute the tests with race detection.")
}
