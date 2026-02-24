package main

import (
	"fmt"
	"strings"
	"time"
)

// 100 Go Mistakes #38, #39: Trim Confusion and Inefficient Concatenation (Solution)
//
// Fixed both issues:
// 1. Use TrimSuffix/TrimPrefix instead of TrimRight/TrimLeft
// 2. Use strings.Builder instead of += for loop concatenation

// cleanFilenames uses TrimSuffix to correctly remove .txt extension.
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
		// Fixed: TrimSuffix removes the exact suffix ".txt".
		cleaned := strings.TrimSuffix(name, ".txt")
		fmt.Printf("  %s -> %s\n", name, cleaned)
	}
}

// cleanPrefixes uses TrimPrefix to correctly remove http:// prefix.
func cleanPrefixes() {
	paths := []string{
		"http://example.com",
		"https://secure.com",
		"http://httpbin.org",
	}

	fmt.Println("Removing http:// prefix:")
	for _, path := range paths {
		// Fixed: TrimPrefix removes the exact prefix "http://".
		cleaned := strings.TrimPrefix(path, "http://")
		fmt.Printf("  %s -> %s\n", path, cleaned)
	}
}

// buildCSV uses strings.Builder for efficient O(n) concatenation.
func buildCSV(rows int) string {
	// Fixed: strings.Builder minimizes allocations by growing an
	// internal buffer. Each WriteString is amortized O(1).
	var b strings.Builder
	for i := 0; i < rows; i++ {
		b.WriteString(fmt.Sprintf("row_%d,value_%d\n", i, i*10))
	}
	return b.String()
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
	test := strings.TrimSuffix("index.html", ".txt")
	expected := "index.html"
	if test == expected {
		fmt.Println("PASS: TrimSuffix correctly leaves non-.txt files alone")
	} else {
		fmt.Printf("FAIL: Got %q, expected %q\n", test, expected)
	}
}
