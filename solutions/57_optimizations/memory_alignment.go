// memory_alignment.go - SOLUTION
// Struct fields reordered for minimal padding, accumulators parallelized.

package main

import (
	"fmt"
	"unsafe"
)

// WastefulStruct has poorly ordered fields (kept for comparison).
type WastefulStruct struct {
	Active    bool
	ID        int64
	Count     int32
	Priority  int16
	Value     float64
	Enabled   bool
	Category  int16
	Flags     int32
	Timestamp int64
	Tag       byte
}

// OptimizedStruct has fields ordered from largest to smallest alignment.
// Fixed: 8-byte fields first, then 4-byte, then 2-byte, then 1-byte.
type OptimizedStruct struct {
	ID        int64   // 8 bytes
	Value     float64 // 8 bytes
	Timestamp int64   // 8 bytes
	Count     int32   // 4 bytes
	Flags     int32   // 4 bytes
	Priority  int16   // 2 bytes
	Category  int16   // 2 bytes
	Active    bool    // 1 byte
	Enabled   bool    // 1 byte
	Tag       byte    // 1 byte + 5 padding to align struct
}

func sumSequential(data []int64) int64 {
	var total int64
	for _, v := range data {
		total += v
	}
	return total
}

// sumParallel uses two accumulators to break the data dependency chain.
// Fixed: even-indexed elements go to total1, odd-indexed to total2.
func sumParallel(data []int64) int64 {
	var total1, total2 int64
	n := len(data)

	// Process pairs of elements into separate accumulators.
	i := 0
	for ; i+1 < n; i += 2 {
		total1 += data[i]
		total2 += data[i+1]
	}

	// Handle the last element if the length is odd.
	if i < n {
		total1 += data[i]
	}

	return total1 + total2
}

func main() {
	fmt.Println("=== Memory Alignment and Data Dependencies ===")

	// Part 1: Struct field ordering
	fmt.Println("\n--- Part 1: Struct Field Ordering ---")

	wastefulSize := unsafe.Sizeof(WastefulStruct{})
	optimizedSize := unsafe.Sizeof(OptimizedStruct{})

	fmt.Printf("WastefulStruct size:   %d bytes\n", wastefulSize)
	fmt.Printf("OptimizedStruct size:  %d bytes\n", optimizedSize)

	saved := int64(wastefulSize) - int64(optimizedSize)
	fmt.Printf("Bytes saved per struct: %d\n", saved)

	if saved > 0 {
		fmt.Println("Struct reordering saved memory!")
		fmt.Printf("For 1 million structs: %d KB saved\n", saved*1_000_000/1_024)
	} else {
		fmt.Println("HINT: Reorder OptimizedStruct fields from largest to smallest alignment")
	}

	// Part 2: Data dependency chains
	fmt.Println("\n--- Part 2: Data Dependency Chains ---")

	data := make([]int64, 10_000_000)
	for i := range data {
		data[i] = int64(i % 100)
	}

	seqResult := sumSequential(data)
	parResult := sumParallel(data)

	fmt.Printf("Sequential sum: %d\n", seqResult)
	fmt.Printf("Parallel sum:   %d\n", parResult)

	if seqResult == parResult {
		fmt.Println("Both produce the same result!")
	} else {
		fmt.Println("ERROR: Results don't match - check your accumulator logic!")
	}

	fmt.Println("\nAlignment optimization complete!")
}
