// race_and_timing_test.go
// Mistakes #83, #86, #87: Fix these tests!
//
// Problems to fix:
// 1. Counter in race_and_timing.go has a data race (Mistake #83)
//    Fix it by adding sync.Mutex to protect concurrent access
// 2. TestConcurrentCounter uses time.Sleep instead of sync.WaitGroup (Mistake #86)
// 3. Cache uses time.Now() directly, making tests slow and flaky (Mistake #87)
//    Refactor Cache to accept a Clock interface, then use a FakeClock in tests
//
// After fixing all issues, delete TestFixRaceAndTiming at the bottom.

package main

import (
	"testing"
	"time"
)

func TestConcurrentCounter(t *testing.T) {
	counter := NewCounter()
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		go func() {
			counter.Increment()
		}()
	}

	// BUG (Mistake #86): Using time.Sleep for synchronization.
	// This is flaky - the goroutines might not finish in time.
	// Replace with sync.WaitGroup to wait for all goroutines.
	time.Sleep(100 * time.Millisecond)

	if counter.Value() != numGoroutines {
		t.Errorf("counter = %d, want %d", counter.Value(), numGoroutines)
	}
}

func TestCacheExpiration(t *testing.T) {
	// BUG (Mistake #87): This test uses a real cache with real time,
	// then sleeps to wait for expiration. Refactor Cache to accept a
	// Clock interface, then use a FakeClock here.
	cache := NewCache(50 * time.Millisecond)
	cache.Set("key", "value")

	val, ok := cache.Get("key")
	if !ok || val != "value" {
		t.Errorf("cache.Get(\"key\") = (%q, %v), want (\"value\", true)", val, ok)
	}

	// BUG: Sleeping to wait for expiration - slow and flaky
	time.Sleep(100 * time.Millisecond)

	_, ok = cache.Get("key")
	if ok {
		t.Error("cache.Get(\"key\") should return false after expiration")
	}
}

func TestCacheOverwrite(t *testing.T) {
	cache := NewCache(1 * time.Minute)
	cache.Set("key", "first")
	cache.Set("key", "second")

	val, ok := cache.Get("key")
	if !ok || val != "second" {
		t.Errorf("cache.Get(\"key\") = (%q, %v), want (\"second\", true)", val, ok)
	}
}

func TestCacheMiss(t *testing.T) {
	cache := NewCache(1 * time.Minute)

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("cache.Get(\"nonexistent\") should return false")
	}
}

// This test forces the exercise to fail until you fix all issues.
// Delete this function after:
// 1. Adding sync.Mutex to Counter (race_and_timing.go)
// 2. Replacing time.Sleep with sync.WaitGroup in TestConcurrentCounter
// 3. Refactoring Cache to accept a Clock interface
// 4. Using a FakeClock in TestCacheExpiration instead of time.Sleep
func TestFixRaceAndTiming(t *testing.T) {
	t.Fatal("EXERCISE: Fix the Counter data race with sync.Mutex, " +
		"replace time.Sleep with sync.WaitGroup, and refactor Cache " +
		"to use a Clock interface for testable time. Then delete this function.")
}
