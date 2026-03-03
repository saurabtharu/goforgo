// variable_shadowing.go
// Understand and fix variable shadowing bugs caused by := in nested scopes
//
// The := operator in Go creates a NEW variable in the current scope, even if
// a variable with the same name exists in an outer scope. This is called
// "variable shadowing" and it's one of the most common Go mistakes.
//
// This program has several shadowing bugs. Each function demonstrates a
// different way that := can silently create a new variable instead of
// modifying the outer one. Fix all the shadowing bugs so the output matches
// the expected results.

package main

import (
	"fmt"
	"strconv"
)

// BUG 1: Shadow in an if block
// The connection string should be updated inside the if block,
// but := creates a new variable that disappears after the if.
func buildConnectionString(host string, port int, useTLS bool) string {
	connStr := host + ":" + strconv.Itoa(port)

	if useTLS {
		// TODO: Fix this shadowing bug. The outer connStr should be modified.
		connStr := connStr + "?tls=true"
		fmt.Println("  TLS enabled:", connStr)
	}

	return connStr
}

// BUG 2: Shadow in a for loop
// The total should accumulate across iterations, but := inside the
// loop body resets it each time.
func sumPositive(numbers []int) int {
	total := 0

	for _, n := range numbers {
		if n > 0 {
			// TODO: Fix this shadowing bug. We want to add to the outer total.
			total := total + n
			_ = total
		}
	}

	return total
}

// BUG 3: Shadow in a switch statement
// The status message should be set inside the switch cases,
// but := creates a new variable in each case block.
func classifyTemperature(celsius float64) string {
	status := "unknown"

	switch {
	case celsius < 0:
		// TODO: Fix this shadowing bug.
		status := "freezing"
		fmt.Println("  Classified as:", status)
	case celsius < 20:
		// TODO: Fix this shadowing bug.
		status := "cold"
		fmt.Println("  Classified as:", status)
	case celsius < 30:
		// TODO: Fix this shadowing bug.
		status := "comfortable"
		fmt.Println("  Classified as:", status)
	default:
		// TODO: Fix this shadowing bug.
		status := "hot"
		fmt.Println("  Classified as:", status)
	}

	return status
}

// BUG 4: Shadow with multiple return values
// When using a function that returns (value, error), := will shadow
// ALL variables on the left side if at least one is new.
func parseConfig(raw string) (int, error) {
	value := 0
	var err error

	if raw != "" {
		// TODO: Fix this shadowing bug. We want to set the outer value and err.
		value, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println("  Parse error:", err)
			return 0, err
		}
		fmt.Println("  Parsed value:", value)
	}

	_ = err
	return value, nil
}

func main() {
	fmt.Println("=== Variable Shadowing Bugs ===")

	fmt.Println("\n--- Bug 1: Shadow in if block ---")
	conn := buildConnectionString("localhost", 5432, true)
	fmt.Println("Connection:", conn)
	// Expected: Connection: localhost:5432?tls=true

	fmt.Println("\n--- Bug 2: Shadow in for loop ---")
	sum := sumPositive([]int{3, -1, 7, -2, 5})
	fmt.Println("Sum of positives:", sum)
	// Expected: Sum of positives: 15

	fmt.Println("\n--- Bug 3: Shadow in switch ---")
	temp := classifyTemperature(22.5)
	fmt.Println("Temperature status:", temp)
	// Expected: Temperature status: comfortable

	fmt.Println("\n--- Bug 4: Shadow with multi-return ---")
	val, err := parseConfig("42")
	fmt.Println("Config value:", val, "Error:", err)
	// Expected: Config value: 42 Error: <nil>

	fmt.Println("\nAll shadowing bugs fixed!")
}
