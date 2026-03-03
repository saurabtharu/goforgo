package main

import (
	"fmt"
	"runtime"
	"strings"
)

// 100 Go Mistakes #40, #41: Useless Conversions and Substring Memory Leaks (Solution)
//
// Fixed both issues:
// 1. Eliminated unnecessary string <-> []byte conversions
// 2. Used strings.Clone to prevent substring memory leaks

// countVowels works entirely with strings — no []byte conversions needed.
func countVowels(s string) int {
	// Fixed: Stay in string-land. Use strings.ToLower and
	// strings.ContainsRune to avoid any []byte allocations.
	lower := strings.ToLower(s)
	count := 0
	for _, ch := range lower {
		if strings.ContainsRune("aeiou", ch) {
			count++
		}
	}
	return count
}

// extractID extracts the ID and copies it to break the backing array reference.
func extractID(logLine string) string {
	idx := strings.Index(logLine, "id=")
	if idx == -1 {
		return ""
	}
	start := idx + 3
	end := strings.IndexByte(logLine[start:], ' ')
	if end == -1 {
		// Fixed: Clone to release reference to the full logLine.
		return strings.Clone(logLine[start:])
	}
	// Fixed: strings.Clone creates an independent copy so the
	// original logLine's backing array can be garbage collected.
	return strings.Clone(logLine[start : start+end])
}

// simulateLogProcessing shows the memory-safe version.
func simulateLogProcessing() {
	padding := strings.Repeat("X", 1024)
	logLine := fmt.Sprintf("timestamp=2024-01-01 level=info %s id=abc123 msg=done", padding)

	fmt.Printf("Original log line size: %d bytes\n", len(logLine))

	ids := make([]string, 100)
	for i := 0; i < 100; i++ {
		ids[i] = extractID(logLine)
	}

	runtime.GC()
	fmt.Printf("Extracted %d IDs, each is: %q\n", len(ids), ids[0])
	fmt.Println("(With strings.Clone, each ID is an independent copy)")

	if ids[0] == "abc123" {
		fmt.Println("PASS: ID extracted correctly")
	} else {
		fmt.Printf("FAIL: Got %q, expected \"abc123\"\n", ids[0])
	}
}

func main() {
	fmt.Println("=== String Conversions and Memory Leaks ===")
	fmt.Println()

	fmt.Println("--- Useless Conversions ---")
	text := "Hello, Beautiful World!"
	count := countVowels(text)
	fmt.Printf("Vowels in %q: %d\n", text, count)
	if count == 8 {
		fmt.Println("PASS: Correct vowel count")
	} else {
		fmt.Println("FAIL: Wrong vowel count")
	}
	fmt.Println()

	fmt.Println("--- Substring Memory Leak ---")
	simulateLogProcessing()
}
