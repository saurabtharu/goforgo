package main

import (
	"encoding/json"
	"fmt"
)

// 100 Go Mistakes #22, #23: Nil vs Empty Slices (Solution)
//
// Fixed:
// 1. Initialize Users as empty slice for correct JSON [] output
// 2. Use len() == 0 to check emptiness (catches both nil and empty)
// 3. Initialize filter result as empty slice for consistent non-nil return

type APIResponse struct {
	Users  []string `json:"users"`
	Errors []string `json:"errors"`
}

func buildResponse(userNames []string) APIResponse {
	// Fixed: Start with an empty slice so JSON marshals to [] not null.
	users := []string{}
	for _, name := range userNames {
		users = append(users, name)
	}

	var errors []string

	return APIResponse{
		Users:  users,
		Errors: errors,
	}
}

func processItems(items []string) {
	// Fixed: Use len() == 0 to catch both nil and empty slices.
	if len(items) == 0 {
		fmt.Println("No items to process")
		return
	}
	fmt.Printf("Processing %d items\n", len(items))
}

func filterPositive(numbers []int) []string {
	// Fixed: Initialize as empty slice so result is always non-nil.
	result := []string{}
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

	fmt.Println("--- JSON Marshaling ---")
	resp := buildResponse(nil)
	data, _ := json.Marshal(resp)
	fmt.Printf("Empty response: %s\n", data)

	resp2 := buildResponse([]string{"Alice", "Bob"})
	data2, _ := json.Marshal(resp2)
	fmt.Printf("With users: %s\n", data2)
	fmt.Println()

	fmt.Println("--- Emptiness Checking ---")
	processItems(nil)
	processItems([]string{})
	processItems([]string{"a"})
	fmt.Println()

	fmt.Println("--- Filter Results ---")
	positives := filterPositive([]int{-1, -2, -3})
	if len(positives) == 0 {
		fmt.Println("No positive numbers found (correct check)")
	} else {
		fmt.Println("Found positives (wrong!)")
	}
}
