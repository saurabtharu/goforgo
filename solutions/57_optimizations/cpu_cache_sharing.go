// cpu_cache_sharing.go - SOLUTION
// Both cache-unfriendly traversal and false sharing are fixed.

package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const matrixSize = 4_096

// sumMatrixColumnMajor iterates in column-major order (cache-unfriendly).
func sumMatrixColumnMajor(matrix *[matrixSize][matrixSize]int64) int64 {
	var total int64
	for col := 0; col < matrixSize; col++ {
		for row := 0; row < matrixSize; row++ {
			total += matrix[row][col]
		}
	}
	return total
}

// sumMatrixRowMajor iterates in row-major order (cache-friendly).
// Fixed: outer loop is rows, inner loop is columns, so memory access
// is sequential and prefetcher-friendly.
func sumMatrixRowMajor(matrix *[matrixSize][matrixSize]int64) int64 {
	var total int64
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			total += matrix[row][col]
		}
	}
	return total
}

// NaiveCounters has two counters sharing a cache line.
type NaiveCounters struct {
	CounterA int64
	CounterB int64
}

// PaddedCounters separates counters onto different cache lines.
// Fixed: 56 bytes of padding after CounterA ensures CounterB starts
// on a new 64-byte cache line boundary.
type PaddedCounters struct {
	CounterA int64
	_pad     [56]byte
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
