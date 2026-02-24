package main

import (
	"fmt"
	"runtime"
)

// 100 Go Mistakes #26, #27, #28: Slice and Map Memory Leaks (Solution)
//
// Fixed:
// 1. Copy sub-slice data to a right-sized slice to free backing array
// 2. Provide size hint to make(map) when size is known
// 3. Demonstrated that maps never shrink and the re-creation pattern

type Message struct {
	ID   int
	Data [1024]byte
}

func getRecentIDs(messages []Message) []int {
	recentCount := 3
	if len(messages) < recentCount {
		recentCount = len(messages)
	}

	ids := make([]int, 0)
	for i := len(messages) - recentCount; i < len(messages); i++ {
		ids = append(ids, messages[i].ID)
	}

	allIDs := make([]int, 0, 1_000)
	for _, m := range messages {
		allIDs = append(allIDs, m.ID)
	}

	// Fixed: Copy sub-slice to a right-sized slice.
	// This lets the large backing array be garbage collected.
	sub := allIDs[len(allIDs)-recentCount:]
	result := make([]int, len(sub))
	copy(result, sub)
	return result
}

func buildLookup(keys []string) map[string]int {
	// Fixed: Provide size hint to avoid repeated rehashing.
	lookup := make(map[string]int, len(keys))

	for i, key := range keys {
		lookup[key] = i
	}

	return lookup
}

func demonstrateMapShrink() {
	m := make(map[int]struct{})

	for i := 0; i < 10_000; i++ {
		m[i] = struct{}{}
	}

	var beforeStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&beforeStats)
	beforeAlloc := beforeStats.HeapAlloc

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
