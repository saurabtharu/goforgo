package main

import "fmt"

// 100 Go Mistakes #24, #25: Slice Copy and Append Pitfalls
//
// This exercise covers two tricky slice mistakes:
// 1. copy() requires the destination to have sufficient length (not just capacity)
// 2. append() can modify the original slice through shared backing array
//
// FIX all the issues so the program works correctly.

// copySlice demonstrates incorrect use of copy().
func copySlice() {
	src := []int{1, 2, 3, 4, 5}

	// FIX: copy() copies min(len(dst), len(src)) elements.
	// A nil/empty destination has length 0, so nothing gets copied!
	// The destination must have sufficient LENGTH, not just capacity.
	dst := make([]int, 0, len(src))
	copy(dst, src)

	fmt.Printf("Source:      %v\n", src)
	fmt.Printf("Destination: %v (len=%d)\n", dst, len(dst))
	// Expected: Destination should have all 5 elements
}

// appendSideEffect shows how append can modify the original.
func appendSideEffect() {
	// Create a base slice and give a "view" to two consumers.
	base := make([]int, 0, 10)
	base = append(base, 1, 2, 3)

	// FIX: consumer1 and consumer2 share base's backing array.
	// Appending to consumer1 overwrites base[3], which consumer2
	// can also see. This is a major source of bugs.
	// Use the full slice expression [0:3:3] to cap the capacity.
	consumer1 := base[0:3]
	consumer2 := base[0:3]

	consumer1 = append(consumer1, 100)
	consumer2 = append(consumer2, 200)

	fmt.Printf("Consumer1: %v\n", consumer1)
	fmt.Printf("Consumer2: %v\n", consumer2)
	// BUG: consumer1's last element gets overwritten by consumer2's append!
	// Expected: consumer1 = [1 2 3 100], consumer2 = [1 2 3 200]
}

// safeSliceFunction demonstrates passing slices to functions safely.
func safeSliceFunction() {
	original := []int{10, 20, 30, 40, 50}

	// FIX: This function receives a slice header (pointer to backing array).
	// If the function appends within capacity, it modifies the original's data.
	// Make an explicit copy before passing to the function.
	modified := addElement(original[:3], 99)

	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Modified: %v\n", modified)
	// Expected: original should be unchanged as [10 20 30 40 50]
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
