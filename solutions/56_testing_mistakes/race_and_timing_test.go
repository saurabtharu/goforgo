// race_and_timing_test.go
// Solution for Mistakes #83, #86, #87
//
// Fixed:
// 1. TestConcurrentCounter uses sync.WaitGroup instead of time.Sleep (Mistake #86)
// 2. TestCacheExpiration uses a FakeClock instead of time.Sleep (Mistake #87)
// 3. Counter itself is fixed with sync.Mutex so -race passes (Mistake #83)

package main

import (
	"sync"
	"testing"
	"time"
)

// FakeClock implements Clock for testing. Lets tests control time
// without sleeping, making them fast and deterministic.
type FakeClock struct {
	now time.Time
}

func NewFakeClock() *FakeClock {
	return &FakeClock{now: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (f *FakeClock) Now() time.Time {
	return f.now
}

func (f *FakeClock) Advance(d time.Duration) {
	f.now = f.now.Add(d)
}

// FIXED: Uses sync.WaitGroup instead of time.Sleep for synchronization.
// This is deterministic - the test waits for exactly all goroutines to finish.
func TestConcurrentCounter(t *testing.T) {
	counter := NewCounter()
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()

	if counter.Value() != numGoroutines {
		t.Errorf("counter = %d, want %d", counter.Value(), numGoroutines)
	}
}

// FIXED: Uses FakeClock to control time without sleeping.
// This test is fast and deterministic.
func TestCacheExpiration(t *testing.T) {
	clock := NewFakeClock()
	cache := NewCache(5*time.Minute, clock)
	cache.Set("key", "value")

	// Verify item exists
	val, ok := cache.Get("key")
	if !ok || val != "value" {
		t.Errorf("cache.Get(\"key\") = (%q, %v), want (\"value\", true)", val, ok)
	}

	// Advance the fake clock past the TTL - no sleeping needed
	clock.Advance(6 * time.Minute)

	// After expiration, item should be gone
	_, ok = cache.Get("key")
	if ok {
		t.Error("cache.Get(\"key\") should return false after expiration")
	}
}

func TestCacheOverwrite(t *testing.T) {
	clock := NewFakeClock()
	cache := NewCache(1*time.Minute, clock)
	cache.Set("key", "first")
	cache.Set("key", "second")

	val, ok := cache.Get("key")
	if !ok || val != "second" {
		t.Errorf("cache.Get(\"key\") = (%q, %v), want (\"second\", true)", val, ok)
	}
}

func TestCacheMiss(t *testing.T) {
	clock := NewFakeClock()
	cache := NewCache(1*time.Minute, clock)

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("cache.Get(\"nonexistent\") should return false")
	}
}
