// cpu_cache_sharing.go
// Understand CPU cache effects and false sharing in concurrent Go programs.
//
// This exercise covers two critical performance mistakes:
//
// 1. Cache-unfriendly access patterns: Iterating over a 2D matrix in
//    column-major order causes constant cache misses because each access
//    jumps to a different cache line. Row-major order is cache-friendly
//    because elements are contiguous in memory.
//
// 2. False sharing: When goroutines update adjacent fields in a struct,
//    those fields likely share a CPU cache line (typically 64 bytes). Each
//    write invalidates the other core's cache, causing massive slowdowns.
//    Padding fields to separate cache lines eliminates this.
//
// Fix both problems so the program prints that the optimized versions
// are faster than the naive ones.

package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const matrixSize = 4_096

// sumMatrixColumnMajor iterates in column-major order (cache-unfriendly).
// Each access jumps matrixSize elements forward in memory.
func sumMatrixColumnMajor(matrix *[matrixSize][matrixSize]int64) int64 {
	var total int64
	for col := 0; col < matrixSize; col++ {
		for row := 0; row < matrixSize; row++ {
			total += matrix[row][col]
		}
	}
	return total
}

// sumMatrixRowMajor iterates in row-major order.
// TODO: Fix this function to iterate in row-major order (cache-friendly).
// Row-major means: iterate rows in the outer loop, columns in the inner loop.
// This accesses memory sequentially, which is friendly to CPU cache prefetching.
func sumMatrixRowMajor(matrix *[matrixSize][matrixSize]int64) int64 {
	var total int64
	// FIX: This is currently also column-major. Change the loop order
	// so the outer loop iterates over rows and the inner loop over columns.
	for col := 0; col < matrixSize; col++ {
		for row := 0; row < matrixSize; row++ {
			total += matrix[row][col]
		}
	}
	return total
}

// NaiveCounters has two counters in adjacent memory locations.
// When two goroutines update these concurrently, they share a cache line,
// causing false sharing and severe performance degradation.
type NaiveCounters struct {
	CounterA int64
	CounterB int64
}

// PaddedCounters separates counters onto different cache lines.
// TODO: Add padding between CounterA and CounterB so they live on
// separate cache lines. A cache line is typically 64 bytes.
// An int64 is 8 bytes, so you need padding of [56]byte (or use [64]byte
// alignment) between the two counters.
type PaddedCounters struct {
	// FIX: Add a padding field after CounterA to push CounterB to a
	// different cache line. Use a [56]byte array as padding.
	CounterA int64
	CounterB int64
}

const iterations = 50_000_000

func benchmarkNaiveCounters() time.Duration {
	counters := &NaiveCounters{}
	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			counters.CounterA++
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			counters.CounterB++
		}
	}()
	wg.Wait()
	return time.Since(start)
}

func benchmarkPaddedCounters() time.Duration {
	counters := &PaddedCounters{}
	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			counters.CounterA++
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			counters.CounterB++
		}
	}()
	wg.Wait()
	return time.Since(start)
}

func main() {
	fmt.Println("=== CPU Cache and False Sharing ===")

	// Part 1: Matrix traversal order
	fmt.Println("\n--- Part 1: Cache-Friendly Matrix Traversal ---")
	var matrix [matrixSize][matrixSize]int64
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			matrix[i][j] = int64(i + j)
		}
	}

	start := time.Now()
	sumCol := sumMatrixColumnMajor(&matrix)
	colTime := time.Since(start)

	start = time.Now()
	sumRow := sumMatrixRowMajor(&matrix)
	rowTime := time.Since(start)

	fmt.Printf("Column-major sum: %d, time: %v\n", sumCol, colTime)
	fmt.Printf("Row-major sum:    %d, time: %v\n", sumRow, rowTime)

	if sumCol != sumRow {
		fmt.Println("ERROR: Sums don't match!")
	}

	if rowTime < colTime {
		fmt.Println("Row-major is faster (cache-friendly access pattern)")
	} else {
		fmt.Println("HINT: Row-major should be faster - fix the iteration order!")
	}

	// Part 2: False sharing
	fmt.Println("\n--- Part 2: False Sharing ---")
	fmt.Printf("NaiveCounters size:  %d bytes\n", unsafe.Sizeof(NaiveCounters{}))
	fmt.Printf("PaddedCounters size: %d bytes\n", unsafe.Sizeof(PaddedCounters{}))

	naiveTime := benchmarkNaiveCounters()
	paddedTime := benchmarkPaddedCounters()

	fmt.Printf("Naive counters:  %v\n", naiveTime)
	fmt.Printf("Padded counters: %v\n", paddedTime)

	if paddedTime < naiveTime {
		fmt.Println("Padded counters are faster (no false sharing)")
	} else {
		fmt.Println("HINT: Padded counters should be faster - add cache line padding!")
	}

	fmt.Println("\nOptimization complete!")
}
