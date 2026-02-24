package main

import (
	"errors"
	"fmt"
)

// 100 Go Mistakes #49: Error Wrapping (Solution)
//
// Fixed all three error wrapping patterns:
// 1. Added context with %w so callers know what failed
// 2. Used %w instead of %v to preserve the error chain
// 3. Used %v (no wrapping) to hide implementation details

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)

// --- Fix 1: Add context with %w ---

func findUser(id int) error {
	if id == 42 {
		return ErrNotFound
	}
	return nil
}

func getUser(id int) error {
	err := findUser(id)
	if err != nil {
		// Fixed: Wrap with context so callers know WHAT was not found.
		// errors.Is(result, ErrNotFound) still works through the chain.
		return fmt.Errorf("getUser id=%d: %w", id, err)
	}
	return nil
}

// --- Fix 2: Use %w to preserve error chain ---

func handleRequest(userID int) error {
	err := getUser(userID)
	if err != nil {
		// Fixed: %w preserves the error chain.
		// errors.Is(result, ErrNotFound) returns true.
		return fmt.Errorf("handleRequest: %w", err)
	}
	return nil
}

// --- Fix 3: Don't wrap internal errors ---

var errRedisDown = errors.New("redis: connection refused")

func connectToRedis() error {
	return errRedisDown
}

func getFromCache(key string) error {
	err := connectToRedis()
	if err != nil {
		// Fixed: Don't wrap the internal error — create a new error
		// that doesn't expose our Redis dependency. Callers can't
		// do errors.Is(err, errRedisDown) anymore.
		return fmt.Errorf("cache unavailable for key %q", key)
	}
	return nil
}

func main() {
	fmt.Println("=== Error Wrapping ===")

	fmt.Println("--- Fix 1: Context Added ---")
	err := getUser(42)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is ErrNotFound: %v\n", errors.Is(err, ErrNotFound))
	fmt.Printf("Has context about what failed: yes\n")

	fmt.Println()

	fmt.Println("--- Fix 2: Correct Verb ---")
	err = handleRequest(42)
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is ErrNotFound: %v (chain preserved by %%w!)\n", errors.Is(err, ErrNotFound))

	fmt.Println()

	fmt.Println("--- Fix 3: Implementation Hidden ---")
	err = getFromCache("session-abc")
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Is errRedisDown: %v (implementation detail hidden!)\n", errors.Is(err, errRedisDown))
}
