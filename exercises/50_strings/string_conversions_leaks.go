package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// 100 Go Mistakes #40, #41: Useless Conversions and Substring Memory Leaks
//
// This exercise covers two string efficiency mistakes:
// 1. Unnecessary string <-> []byte conversions that waste allocations
// 2. Substring memory leaks where a small substring keeps the entire
//    original string's backing array alive in memory
//
// FIX both issues so the program runs correctly and efficiently.

// countVowels demonstrates useless string/[]byte conversions.
// It converts back and forth between string and []byte unnecessarily.
func countVowels(s string) int {
	// BUG: This function converts string -> []byte -> string -> []byte
	// repeatedly. Each conversion allocates a new copy. The bytes package
	// can work directly with []byte, and strings package works with strings.
	// FIX: Pick one type and stick with it. Since we receive a string,
	// use the strings package functions directly — no []byte needed.
	b := []byte(s)                        // unnecessary conversion 1
	lower := strings.ToLower(string(b))   // unnecessary conversion 2 (back to string)
	data := []byte(lower)                 // unnecessary conversion 3

	count := 0
	vowels := []byte("aeiou")
	for _, ch := range data {
		if bytes.Contains(vowels, []byte{ch}) { // unnecessary: could use strings.ContainsRune
			count++
		}
	}
	return count
}

// extractID simulates extracting a small substring from a large log line.
// This demonstrates the substring memory leak pattern.
func extractID(logLine string) string {
	// BUG: In Go, a substring operation like s[start:end] shares the
	// same backing array as the original string. If logLine is a 10KB
	// log entry and we extract a 10-byte ID, the 10-byte substring
	// still references the entire 10KB backing array, preventing it
	// from being garbage collected.
	//
	// FIX: Copy the substring to a new string to release the reference
	// to the original backing array. Use strings.Clone() (Go 1.20+)
	// or string([]byte(sub)) for older versions.
	idx := strings.Index(logLine, "id=")
	if idx == -1 {
		return ""
	}
	start := idx + 3
	end := strings.IndexByte(logLine[start:], ' ')
	if end == -1 {
		return logLine[start:]
	}
	return logLine[start : start+end]
}

// simulateLogProcessing shows the memory impact of the substring leak.
func simulateLogProcessing() {
	// Build a large log line (~1KB of padding + small ID)
	padding := strings.Repeat("X", 1024)
	logLine := fmt.Sprintf("timestamp=2024-01-01 level=info %s id=abc123 msg=done", padding)

	fmt.Printf("Original log line size: %d bytes\n", len(logLine))

	// Extract many IDs, keeping references
	ids := make([]string, 100)
	for i := 0; i < 100; i++ {
		ids[i] = extractID(logLine)
	}

	runtime.GC()
	fmt.Printf("Extracted %d IDs, each is: %q\n", len(ids), ids[0])
	fmt.Println("(Without copying, each small ID holds a reference to the full log line)")

	// Verify the ID was extracted correctly
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
