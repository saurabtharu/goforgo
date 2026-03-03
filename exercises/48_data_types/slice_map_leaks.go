package main

import (
	"fmt"
	"runtime"
)

// 100 Go Mistakes #26, #27, #28: Slice and Map Memory Leaks
//
// This exercise covers three memory-related mistakes:
// 1. Slice memory leaks: sub-slice keeps entire backing array alive
// 2. Inefficient map initialization: not providing size hints
// 3. Maps never shrink: deleting keys doesn't free memory
//
// FIX all three issues to demonstrate proper memory management.

// Message simulates a large struct stored in a slice.
type Message struct {
	ID   int
	Data [1024]byte // 1KB payload per message
}

// getRecentIDs demonstrates the slice memory leak pattern.
// It receives a large slice but only needs a few element IDs.
func getRecentIDs(messages []Message) []int {
	// FIX: Taking a sub-slice of IDs derived from a range is fine,
	// but building a result slice from the messages' data and then
	// returning a sub-view of the original keeps the entire
	// backing array alive in memory.
	//
	// The correct fix is to copy the needed elements into a new slice
	// so the large backing array can be garbage collected.
	recentCount := 3
	if len(messages) < recentCount {
		recentCount = len(messages)
	}

	// This creates IDs referencing only what we need - looks correct but
	// demonstrates the pattern. The real leak happens with slice-of-struct.
	ids := make([]int, 0)
	for i := len(messages) - recentCount; i < len(messages); i++ {
		ids = append(ids, messages[i].ID)
	}

	// BUG: Now imagine we stored the messages sub-slice instead:
	// recentMsgs := messages[len(messages)-recentCount:]
	// This would keep ALL messages alive! Instead, we must copy.
	// Demonstrate by returning a sub-slice that retains large capacity:
	allIDs := make([]int, 0, 1_000)
	for _, m := range messages {
		allIDs = append(allIDs, m.ID)
	}

	// FIX: This sub-slice retains backing array of capacity 1000.
	// Copy to a right-sized slice instead.
	result := allIDs[len(allIDs)-recentCount:]
	return result
}

// buildLookup demonstrates inefficient map initialization.
func buildLookup(keys []string) map[string]int {
	// FIX: Not providing a size hint causes the map to rehash
	// multiple times as it grows. When you know the approximate
	// size, pass it to make().
	lookup := make(map[string]int)

	for i, key := range keys {
		lookup[key] = i
	}

	return lookup
}

// demonstrateMapShrink shows that maps never shrink after deletions.
func demonstrateMapShrink() {
	m := make(map[int]struct{})

	// Fill the map
	for i := 0; i < 10_000; i++ {
		m[i] = struct{}{}
	}

	var beforeStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&beforeStats)
	beforeAlloc := beforeStats.HeapAlloc

	// Delete all entries
	for i := 0; i < 10_000; i++ {
		delete(m, i)
	}

	var afterStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&afterStats)
	afterAlloc := afterStats.HeapAlloc

	fmt.Printf("Map length after deletion: %d\n", len(m))
	fmt.Printf("Heap before deletion: %d KB\n", beforeAlloc/1024)
	fmt.Printf("Heap after deletion:  %d KB\n", afterAlloc/1024)
	fmt.Println("Note: Map memory is NOT fully reclaimed by deletion alone.")

	// FIX: To truly free map memory, create a new map and copy
	// any remaining entries. Show the correct pattern:
	fmt.Println()
	fmt.Println("To reclaim map memory, create a new map:")
	fmt.Println("  newMap := make(map[int]struct{}, len(oldMap))")
	fmt.Println("  for k, v := range oldMap { newMap[k] = v }")
}

func main() {
	fmt.Println("=== Slice and Map Memory Leaks ===")
	fmt.Println()

	fmt.Println("--- Slice Memory Leak ---")
	messages := make([]Message, 100)
	for i := range messages {
		messages[i].ID = i + 1
	}
	recent := getRecentIDs(messages)
	fmt.Printf("Recent IDs: %v\n", recent)
	fmt.Printf("Result length: %d, capacity: %d\n", len(recent), cap(recent))
	if cap(recent) == len(recent) {
		fmt.Println("PASS: No wasted capacity (memory leak fixed)")
	} else {
		fmt.Println("FAIL: Excess capacity retained (memory leak!)")
	}
	fmt.Println()

	fmt.Println("--- Map Size Hint ---")
	keys := make([]string, 1_000)
	for i := range keys {
		keys[i] = fmt.Sprintf("key_%d", i)
	}
	lookup := buildLookup(keys)
	fmt.Printf("Map size: %d\n", len(lookup))
	fmt.Println()

	fmt.Println("--- Maps Never Shrink ---")
	demonstrateMapShrink()
}
