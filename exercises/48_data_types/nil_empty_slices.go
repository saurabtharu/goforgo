package main

import (
	"encoding/json"
	"fmt"
)

// 100 Go Mistakes #22, #23: Nil vs Empty Slices
//
// This exercise covers two related mistakes:
// 1. Nil slice vs empty slice have different semantics (especially in JSON)
// 2. Checking slice emptiness wrong (== nil vs len() == 0)
//
// FIX the issues so JSON output is correct and emptiness checks work properly.

type APIResponse struct {
	Users  []string `json:"users"`
	Errors []string `json:"errors"`
}

// buildResponse creates an API response. The developer wants:
// - "users" to be an empty JSON array [] when there are no users (not null)
// - "errors" to be omitted or null when there are no errors
func buildResponse(userNames []string) APIResponse {
	// FIX: A nil slice marshals to "null" in JSON.
	// An empty slice ([]string{}) marshals to "[]".
	// The developer wants Users to be [] not null when empty.
	var users []string
	for _, name := range userNames {
		users = append(users, name)
	}

	// The errors field should be nil (marshals to "null") when no errors exist.
	var errors []string

	return APIResponse{
		Users:  users,
		Errors: errors,
	}
}

// processItems demonstrates the wrong way to check for empty slices.
func processItems(items []string) {
	// FIX: Checking items == nil misses the case where items is
	// an empty slice ([]string{}). Both nil and empty slices have
	// len() == 0, but only nil slices are == nil.
	// Use len(items) == 0 to catch both cases.
	if items == nil {
		fmt.Println("No items to process")
		return
	}
	fmt.Printf("Processing %d items\n", len(items))
}

// filterPositive shows how append with nil starting slice works fine,
// but the caller might check the result with == nil incorrectly.
func filterPositive(numbers []int) []string {
	// FIX: If no positive numbers exist, result stays nil.
	// The caller checks with == nil and gets confused.
	// Initialize as empty slice so the result is always non-nil.
	var result []string
	for _, n := range numbers {
		if n > 0 {
			result = append(result, fmt.Sprintf("%d", n))
		}
	}
	return result
}

func main() {
	fmt.Println("=== Nil vs Empty Slices ===")
	fmt.Println()

	// Test 1: JSON marshaling with no users
	fmt.Println("--- JSON Marshaling ---")
	resp := buildResponse(nil)
	data, _ := json.Marshal(resp)
	fmt.Printf("Empty response: %s\n", data)
	// Expected: {"users":[],"errors":null}
	// Buggy:   {"users":null,"errors":null}

	// With users
	resp2 := buildResponse([]string{"Alice", "Bob"})
	data2, _ := json.Marshal(resp2)
	fmt.Printf("With users: %s\n", data2)
	fmt.Println()

	// Test 2: Emptiness checking
	fmt.Println("--- Emptiness Checking ---")
	processItems(nil)
	processItems([]string{})     // BUG: this passes the nil check but is still empty!
	processItems([]string{"a"})
	fmt.Println()

	// Test 3: Filter results
	fmt.Println("--- Filter Results ---")
	positives := filterPositive([]int{-1, -2, -3})
	if len(positives) == 0 {
		fmt.Println("No positive numbers found (correct check)")
	} else {
		fmt.Println("Found positives (wrong!)")
	}
}
