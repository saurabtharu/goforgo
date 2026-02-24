package main

import "fmt"

// 100 Go Mistakes #24, #25: Slice Copy and Append Pitfalls (Solution)
//
// Fixed:
// 1. Set destination length correctly for copy()
// 2. Use full slice expression [0:3:3] to prevent append side effects
// 3. Use explicit copy when passing slices to functions

func copySlice() {
	src := []int{1, 2, 3, 4, 5}

	// Fixed: Set destination length equal to source length.
	// copy() copies min(len(dst), len(src)) elements.
	dst := make([]int, len(src))
	copy(dst, src)

	fmt.Printf("Source:      %v\n", src)
	fmt.Printf("Destination: %v (len=%d)\n", dst, len(dst))
}

func appendSideEffect() {
	base := make([]int, 0, 10)
	base = append(base, 1, 2, 3)

	// Fixed: Use full slice expression [0:3:3] to cap capacity at length.
	// This forces append to allocate a new backing array for each consumer.
	consumer1 := base[0:3:3]
	consumer2 := base[0:3:3]

	consumer1 = append(consumer1, 100)
	consumer2 = append(consumer2, 200)

	fmt.Printf("Consumer1: %v\n", consumer1)
	fmt.Printf("Consumer2: %v\n", consumer2)
}

func safeSliceFunction() {
	original := []int{10, 20, 30, 40, 50}

	// Fixed: Use full slice expression to cap capacity, preventing
	// addElement's append from overwriting original's data.
	modified := addElement(original[:3:3], 99)

	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Modified: %v\n", modified)
}

func addElement(s []int, elem int) []int {
	return append(s, elem)
}

func main() {
	fmt.Println("=== Slice Copy and Append Pitfalls ===")
	fmt.Println()

	fmt.Println("--- Incorrect copy() ---")
	copySlice()
	fmt.Println()

	fmt.Println("--- Append Side Effects ---")
	appendSideEffect()
	fmt.Println()

	fmt.Println("--- Safe Slice Passing ---")
	safeSliceFunction()
}
