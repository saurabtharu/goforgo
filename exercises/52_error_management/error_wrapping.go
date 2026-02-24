package main

import (
	"errors"
	"fmt"
)

// 100 Go Mistakes #49: Error Wrapping
//
// Three bugs to fix:
// 1. getUser() returns bare errors with no context — caller can't tell WHAT failed
// 2. handleRequest() uses %v instead of %w, breaking the error chain
// 3. getFromCache() wraps an internal error with %w, leaking implementation details
//
// FIX all three patterns to use proper error wrapping.

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)

// --- Bug 1: No context added to errors ---

func findUser(id int) error {
	if id == 42 {
		return ErrNotFound
	}
	return nil
}

func getUser(id int) error {
	err := findUser(id)
	if err != nil {
		// BUG: Returns the bare error with no context.
		// The caller has no idea WHAT was "not found" — a user? a file? a record?
		return err
	}
	return nil
}

// --- Bug 2: %v instead of %w ---

func handleRequest(userID int) error {
	err := getUser(userID)
	if err != nil {
		// BUG: %v converts the error to a string, creating a NEW error.
		// The original error chain is lost — errors.Is(result, ErrNotFound)
		// will return false because ErrNotFound is no longer in the chain.
		return fmt.Errorf("request failed: %v", err)
	}
	return nil
}

// --- Bug 3: Wrapping leaks implementation details ---

var errRedisDown = errors.New("redis: connection refused")

func connectToRedis() error {
	return errRedisDown
}

func getFromCache(key string) error {
	err := connectToRedis()
	if err != nil {
		// BUG: Wrapping with %w exposes that we use Redis internally.
		// Callers can now do errors.Is(err, errRedisDown) and couple
		// to our implementation. If we switch to memcached, callers break.
		return fmt.Errorf("cache error for %q: %w", key, err)
	}
	return nil
}

func main() {
	fmt.Println("=== Error Wrapping ===")

	fmt.Println("--- Bug 1: No Context ---")
	err := getUser(42)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is ErrNotFound: %v\n", errors.Is(err, ErrNotFound))
	fmt.Printf("Has context about what failed: no\n")

	fmt.Println()

	fmt.Println("--- Bug 2: Wrong Verb ---")
	err = handleRequest(42)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is ErrNotFound: %v (chain broken by %%v!)\n", errors.Is(err, ErrNotFound))

	fmt.Println()

	fmt.Println("--- Bug 3: Leaking Implementation ---")
	err = getFromCache("session-abc")
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is errRedisDown: %v (leaks Redis dependency!)\n", errors.Is(err, errRedisDown))
}
