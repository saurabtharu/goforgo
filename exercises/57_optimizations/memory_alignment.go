// memory_alignment.go
// Understand how struct field ordering affects memory usage in Go.
//
// This exercise covers two important optimization concepts:
//
// 1. Struct padding waste: The Go compiler adds padding bytes between fields
//    to satisfy alignment requirements. A bool (1 byte) followed by an int64
//    (8-byte aligned) wastes 7 bytes of padding. Reordering fields from
//    largest to smallest minimizes this waste.
//
// 2. Data dependency awareness: When processing data, independent operations
//    can be structured to help the CPU pipeline execute more efficiently.
//    Accumulating into separate variables and combining at the end reduces
//    data dependency chains.
//
// Fix the struct field ordering and the accumulator pattern so the
// program reports improved memory usage and processing.

package main

import (
	"fmt"
	"unsafe"
)

// WastefulStruct has poorly ordered fields that waste memory due to padding.
// The compiler must add padding to align each field to its natural boundary:
//   - bool: 1 byte aligned
//   - int64: 8 byte aligned
//   - int32: 4 byte aligned
//   - int16: 2 byte aligned
//   - float64: 8 byte aligned
//
// TODO: Reorder the fields from largest alignment to smallest to minimize padding.
// Group 8-byte fields together, then 4-byte, then 2-byte, then 1-byte.
type WastefulStruct struct {
	Active    bool    // 1 byte + 7 padding
	ID        int64   // 8 bytes
	Count     int32   // 4 bytes + 4 padding
	Priority  int16   // 2 bytes + 6 padding
	Value     float64 // 8 bytes
	Enabled   bool    // 1 byte + 1 padding
	Category  int16   // 2 bytes + 4 padding
	Flags     int32   // 4 bytes
	Timestamp int64   // 8 bytes
	Tag       byte    // 1 byte + 7 padding
}

// OptimizedStruct should have the same fields but reordered to minimize padding.
// FIX: Reorder fields from largest to smallest alignment requirement.
// Put all 8-byte fields (int64, float64) first, then 4-byte (int32),
// then 2-byte (int16), then 1-byte (bool, byte).
type OptimizedStruct struct {
	Active    bool    // 1 byte + 7 padding
	ID        int64   // 8 bytes
	Count     int32   // 4 bytes + 4 padding
	Priority  int16   // 2 bytes + 6 padding
	Value     float64 // 8 bytes
	Enabled   bool    // 1 byte + 1 padding
	Category  int16   // 2 bytes + 4 padding
	Flags     int32   // 4 bytes
	Timestamp int64   // 8 bytes
	Tag       byte    // 1 byte + 7 padding
}

// sumSequential has a long dependency chain: each addition depends on the
// previous result, preventing the CPU from pipelining operations.
func sumSequential(data []int64) int64 {
	var total int64
	for _, v := range data {
		total += v // each iteration depends on the previous total
	}
	return total
}

// sumParallel should use multiple accumulators to break the dependency chain.
// TODO: Use two (or more) separate accumulator variables, add alternating
// elements to each, then combine them at the end. This allows the CPU to
// execute additions in parallel since the accumulators are independent.
func sumParallel(data []int64) int64 {
	// FIX: Use two accumulators (total1, total2) and add even-indexed
	// elements to total1 and odd-indexed elements to total2.
	// Then return total1 + total2.
	var total int64
	for _, v := range data {
		total += v
	}
	return total
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
