package main

import "fmt"

// 100 Go Mistakes #20, #21: Slice Length vs Capacity (Solution)
//
// Fixed all slice length/capacity issues:
// 1. Use make([]int, 0, 3) to pre-allocate capacity without setting length
// 2. Pre-allocate with make([]int, 0, size) to avoid reallocations
// 3. Use full slice expression [:2:2] to limit capacity of subslice

func demonstrateLenCap() {
	// Fixed: make([]int, 0, 3) sets length=0, capacity=3.
	// Now append fills from index 0 as expected.
	s := make([]int, 0, 3)
	s = append(s, 1, 2, 3)

	fmt.Printf("Slice: %v\n", s)
	fmt.Printf("Length: %d, Capacity: %d\n", len(s), cap(s))
}

func buildSliceInefficient() []int {
	size := 1_000

	// Fixed: Pre-allocate capacity to avoid repeated reallocations.
	result := make([]int, 0, size)

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

func appendToSubslice() {
	original := []int{1, 2, 3, 4, 5}

	// Fixed: Use full slice expression [:2:2] to cap the capacity.
	// This forces append to allocate a new backing array.
	sub := original[:2:2]
	sub = append(sub, 99)

	fmt.Printf("Original after append to sub: %v\n", original)
	fmt.Printf("Sub: %v\n", sub)
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
