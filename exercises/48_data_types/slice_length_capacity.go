package main

import "fmt"

// 100 Go Mistakes #20, #21: Slice Length vs Capacity
//
// This exercise covers two common slice mistakes:
// 1. Confusing slice length and capacity
// 2. Inefficient slice initialization (growing vs pre-allocating)
//
// FIX all the issues so the program prints correct results.

// demonstrateLenCap shows the confusion between length and capacity.
func demonstrateLenCap() {
	// FIX: The developer wants a slice with 3 pre-allocated slots to fill in.
	// But make([]int, 3) creates a slice with length=3 AND capacity=3,
	// meaning it already has 3 zero-valued elements.
	// Appending adds BEYOND those 3 zeros.
	s := make([]int, 3)
	s = append(s, 1, 2, 3)

	fmt.Printf("Slice: %v\n", s)
	fmt.Printf("Length: %d, Capacity: %d\n", len(s), cap(s))

	// Expected: [1, 2, 3] with length 3
}

// buildSliceInefficient grows a slice one element at a time without
// pre-allocating, causing unnecessary memory allocations.
func buildSliceInefficient() []int {
	size := 1_000

	// FIX: This creates an empty slice and grows it repeatedly.
	// Each time the backing array fills up, Go allocates a new,
	// larger array and copies everything over.
	// Pre-allocate with make([]int, 0, size) to avoid this.
	var result []int

	allocsBefore := cap(result)
	reallocCount := 0

	for i := 0; i < size; i++ {
		result = append(result, i*2)
		if cap(result) != allocsBefore {
			reallocCount++
			allocsBefore = cap(result)
		}
	}

	fmt.Printf("Built slice of %d elements\n", len(result))
	fmt.Printf("Reallocations: %d\n", reallocCount)

	return result
}

// appendToSubslice shows how length affects append behavior.
func appendToSubslice() {
	original := []int{1, 2, 3, 4, 5}

	// FIX: Taking a subslice s[:2] gives length=2 but capacity=5.
	// Appending to it will overwrite original[2] because there's
	// room in the backing array.
	// The developer expects sub to be independent of original.
	sub := original[:2]
	sub = append(sub, 99)

	fmt.Printf("Original after append to sub: %v\n", original)
	fmt.Printf("Sub: %v\n", sub)

	// Expected: original should remain [1, 2, 3, 4, 5]
}

func main() {
	fmt.Println("=== Slice Length vs Capacity ===")
	fmt.Println()

	fmt.Println("--- Length vs Capacity Confusion ---")
	demonstrateLenCap()
	fmt.Println()

	fmt.Println("--- Inefficient Slice Growth ---")
	buildSliceInefficient()
	fmt.Println()

	fmt.Println("--- Subslice Append Surprise ---")
	appendToSubslice()
}
