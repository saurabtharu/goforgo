package main

import (
	"fmt"
	"unicode/utf8"
)

// 100 Go Mistakes #36, #37: Runes and Iteration (Solution)
//
// Fixed both issues:
// 1. Use range loop for proper rune iteration
// 2. Use utf8.RuneCountInString and []rune conversion for safe operations

// printCharacters prints each rune of a string correctly using range.
func printCharacters(s string) {
	fmt.Printf("String: %q\n", s)

	// Fixed: range loop decodes UTF-8 and yields (byte_index, rune) pairs.
	for i, r := range s {
		fmt.Printf("  index %d: %c\n", i, r)
	}
}

// truncateToN returns the first n runes of a string safely.
func truncateToN(s string, n int) string {
	// Fixed: Convert to []rune for character-based slicing.
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

func main() {
	fmt.Println("=== Runes and Iteration ===")
	fmt.Println()

	greeting := "Hello, 世界"
	fmt.Println("--- Printing Characters ---")
	printCharacters(greeting)
	fmt.Println()

	fmt.Println("--- Truncating Strings ---")
	truncated := truncateToN(greeting, 8)
	fmt.Printf("First 8 characters of %q: %q\n", greeting, truncated)

	expected := "Hello, 世"
	if truncated == expected {
		fmt.Println("PASS: Correct character-based truncation")
	} else {
		fmt.Println("FAIL: Truncation is wrong (probably sliced at byte boundary)")
	}
	fmt.Println()

	emoji := "Go🚀Fun"
	fmt.Printf("String: %q\n", emoji)
	fmt.Printf("  len() reports: %d (this is the byte count)\n", len(emoji))
	// Fixed: Use utf8.RuneCountInString for accurate rune count.
	fmt.Printf("  Rune count: %d\n", utf8.RuneCountInString(emoji))
	runeCount := utf8.RuneCountInString(emoji)
	if runeCount == 6 {
		fmt.Println("PASS: Correct rune count")
	} else {
		fmt.Println("FAIL: Rune count is wrong")
	}
}
