package main

import (
	"fmt"
)

// 100 Go Mistakes #36, #37: Runes and Iteration
//
// This exercise covers two common string mistakes in Go:
// 1. Not understanding that a string is a sequence of bytes, not runes
// 2. Confusing len() (byte count) with rune count
//
// FIX both functions so the program prints correct results.

// printCharacters attempts to print each character of a string
// containing multi-byte Unicode characters.
func printCharacters(s string) {
	fmt.Printf("String: %q\n", s)

	// BUG: Iterating by byte index over a string with multi-byte
	// characters produces garbled output. Each Chinese character
	// is 3 bytes in UTF-8, so s[i] gives individual bytes, not runes.
	// FIX: Use a range loop to iterate over runes instead of bytes.
	for i := 0; i < len(s); i++ {
		fmt.Printf("  index %d: %c\n", i, s[i])
	}
}

// truncateToN attempts to return the first n "characters" of a string.
func truncateToN(s string, n int) string {
	// BUG: len(s) returns the byte count, not the rune (character) count.
	// Slicing s[:n] at a byte offset can cut a multi-byte character in half,
	// producing invalid UTF-8 or wrong output.
	// FIX: Use utf8.RuneCountInString for length, and convert to []rune
	// for safe character-based slicing.
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func main() {
	fmt.Println("=== Runes and Iteration ===")
	fmt.Println()

	// "Hello, 世界" contains ASCII + 3-byte Chinese characters
	greeting := "Hello, 世界"
	fmt.Println("--- Printing Characters ---")
	printCharacters(greeting)
	fmt.Println()

	fmt.Println("--- Truncating Strings ---")
	// We want the first 8 characters: "Hello, 世界" has 9 runes
	truncated := truncateToN(greeting, 8)
	fmt.Printf("First 8 characters of %q: %q\n", greeting, truncated)

	// Verify: the truncated string should be "Hello, 世"
	expected := "Hello, 世"
	if truncated == expected {
		fmt.Println("PASS: Correct character-based truncation")
	} else {
		fmt.Println("FAIL: Truncation is wrong (probably sliced at byte boundary)")
	}
	fmt.Println()

	// Demonstrate rune count vs byte count
	emoji := "Go🚀Fun"
	fmt.Printf("String: %q\n", emoji)
	fmt.Printf("  len() reports: %d (this is the byte count)\n", len(emoji))
	// FIX: Print the correct rune count using utf8.RuneCountInString
	// instead of len(). You will need to add "unicode/utf8" to the imports.
	fmt.Printf("  Rune count: %d\n", len(emoji))
	runeCount := 9 // FIX: this is the byte count, not rune count!
	if runeCount == 6 {
		fmt.Println("PASS: Correct rune count")
	} else {
		fmt.Println("FAIL: Rune count is wrong")
	}
}
