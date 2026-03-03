package main

import (
	"fmt"
	"strings"
	"time"
)

// 100 Go Mistakes #38, #39: Trim Confusion and Inefficient Concatenation
//
// This exercise covers two common string mistakes:
// 1. Confusing TrimRight/TrimLeft (character set removal) with
//    TrimSuffix/TrimPrefix (exact substring removal)
// 2. Using += for string concatenation in loops (O(n^2) allocations)
//
// FIX both issues so the program prints correct results.

// cleanFilenames demonstrates the TrimRight vs TrimSuffix confusion.
func cleanFilenames() {
	filenames := []string{
		"report.txt",
		"photo.png",
		"notes.txt.bak",
		"data.csv",
		"index.html",
	}

	fmt.Println("Removing .txt extension:")
	for _, name := range filenames {
		// BUG: TrimRight removes all trailing characters that appear
		// in the cutset ".txt". This means it strips any trailing
		// combination of '.', 't', 'x' — not just the suffix ".txt".
		// For "index.html" it would strip the trailing 't' from "xt" match.
		// FIX: Use strings.TrimSuffix to remove the exact suffix ".txt".
		cleaned := strings.TrimRight(name, ".txt")
		fmt.Printf("  %s -> %s\n", name, cleaned)
	}
}

// cleanPrefixes demonstrates TrimLeft vs TrimPrefix confusion.
func cleanPrefixes() {
	paths := []string{
		"http://example.com",
		"https://secure.com",
		"http://httpbin.org",
	}

	fmt.Println("Removing http:// prefix:")
	for _, path := range paths {
		// BUG: TrimLeft removes all leading characters in the cutset.
		// For "https://secure.com", it strips 'h','t','p' individually,
		// eating into "https" and producing garbled output.
		// FIX: Use strings.TrimPrefix to remove the exact prefix "http://".
		cleaned := strings.TrimLeft(path, "http://")
		fmt.Printf("  %s -> %s\n", path, cleaned)
	}
}

// buildCSV shows inefficient string concatenation with +=.
func buildCSV(rows int) string {
	// BUG: Using += in a loop creates a new string allocation on every
	// iteration. For n rows this is O(n^2) because each += copies the
	// entire accumulated string.
	// FIX: Use strings.Builder for O(n) concatenation.
	result := ""
	for i := 0; i < rows; i++ {
		result += fmt.Sprintf("row_%d,value_%d\n", i, i*10)
	}
	return result
}

func main() {
	fmt.Println("=== Trim and Concatenation ===")
	fmt.Println()

	fmt.Println("--- TrimRight vs TrimSuffix ---")
	cleanFilenames()
	fmt.Println()

	fmt.Println("--- TrimLeft vs TrimPrefix ---")
	cleanPrefixes()
	fmt.Println()

	fmt.Println("--- String Concatenation Performance ---")
	rows := 10_000

	start := time.Now()
	csv := buildCSV(rows)
	elapsed := time.Since(start)
	_ = csv

	fmt.Printf("Built CSV with %d rows in %v\n", rows, elapsed)
	fmt.Println("(With strings.Builder this should be significantly faster)")

	// Verify TrimSuffix behavior
	fmt.Println()
	test := strings.TrimRight("index.html", ".txt")
	expected := "index.html" // .html should be untouched
	if test == expected {
		fmt.Println("PASS: TrimSuffix correctly leaves non-.txt files alone")
	} else {
		fmt.Printf("FAIL: Got %q, expected %q\n", test, expected)
	}
}
