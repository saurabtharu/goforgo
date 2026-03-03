// variable_shadowing.go - SOLUTION
// All shadowing bugs fixed by using = instead of := in inner scopes.

package main

import (
	"fmt"
	"strconv"
)

// FIX 1: Use = instead of := in the if block.
func buildConnectionString(host string, port int, useTLS bool) string {
	connStr := host + ":" + strconv.Itoa(port)

	if useTLS {
		// Fixed: = instead of := so we modify the outer connStr.
		connStr = connStr + "?tls=true"
		fmt.Println("  TLS enabled:", connStr)
	}

	return connStr
}

// FIX 2: Use = instead of := in the for loop body.
func sumPositive(numbers []int) int {
	total := 0

	for _, n := range numbers {
		if n > 0 {
			// Fixed: = instead of := so we accumulate into the outer total.
			total = total + n
		}
	}

	return total
}

// FIX 3: Use = instead of := in switch cases.
func classifyTemperature(celsius float64) string {
	status := "unknown"

	switch {
	case celsius < 0:
		// Fixed: = instead of :=
		status = "freezing"
		fmt.Println("  Classified as:", status)
	case celsius < 20:
		status = "cold"
		fmt.Println("  Classified as:", status)
	case celsius < 30:
		status = "comfortable"
		fmt.Println("  Classified as:", status)
	default:
		status = "hot"
		fmt.Println("  Classified as:", status)
	}

	return status
}

// FIX 4: Separate the new variable (err) from the assignment.
func parseConfig(raw string) (int, error) {
	value := 0
	var err error

	if raw != "" {
		// Fixed: use = for value, declare err in the outer scope already.
		value, err = strconv.Atoi(raw)
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

	fmt.Println("\n--- Bug 2: Shadow in for loop ---")
	sum := sumPositive([]int{3, -1, 7, -2, 5})
	fmt.Println("Sum of positives:", sum)

	fmt.Println("\n--- Bug 3: Shadow in switch ---")
	temp := classifyTemperature(22.5)
	fmt.Println("Temperature status:", temp)

	fmt.Println("\n--- Bug 4: Shadow with multi-return ---")
	val, err := parseConfig("42")
	fmt.Println("Config value:", val, "Error:", err)

	fmt.Println("\nAll shadowing bugs fixed!")
}
