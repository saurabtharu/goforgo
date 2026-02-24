// concurrency_vs_parallelism.go
// Understand the difference between concurrency and parallelism,
// and learn that concurrency is NOT always faster.
// Based on 100 Go Mistakes #55 and #56.
//
// Concurrency is about STRUCTURE: dealing with multiple things at once.
// Parallelism is about EXECUTION: doing multiple things at once.
// Spawning goroutines for tiny work units wastes more on overhead
// than you gain from parallel execution.
//
// I AM NOT DONE YET!

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var goroutinesSpawned int64

// merge combines two sorted slices into one sorted slice.
func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// sequentialMergeSort sorts without goroutines.
func sequentialMergeSort(data []int) []int {
	if len(data) <= 1 {
		return data
	}
	mid := len(data) / 2
	left := sequentialMergeSort(data[:mid])
	right := sequentialMergeSort(data[mid:])
	return merge(left, right)
}

// concurrentMergeSort sorts using goroutines for parallelism.
//
// FIX: This function spawns goroutines for EVERY recursive call, even for
// slices of 2-3 elements. Each goroutine costs ~2KB of stack space plus
// scheduler overhead, which far exceeds the cost of sorting a tiny slice.
//
// Add a const `concurrencyThreshold = 4`. When len(data) is at or below the
// threshold, sort sequentially instead of spawning new goroutines.
func concurrentMergeSort(data []int) []int {
	if len(data) <= 1 {
		return data
	}
	mid := len(data) / 2

	// BUG: Always spawns goroutines regardless of slice size.
	// For a 16-element array this creates 30 goroutines!
	var left, right []int
	var wg sync.WaitGroup
	wg.Add(2)

	atomic.AddInt64(&goroutinesSpawned, 2)
	go func() {
		defer wg.Done()
		left = concurrentMergeSort(data[:mid])
	}()
	go func() {
		defer wg.Done()
		right = concurrentMergeSort(data[mid:])
	}()
	wg.Wait()

	return merge(left, right)
}

func isSorted(data []int) bool {
	for i := 1; i < len(data); i++ {
		if data[i] < data[i-1] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("=== Mistake #55: Concurrency vs Parallelism ===")
	fmt.Println("Concurrency = structure (dealing with multiple things)")
	fmt.Println("Parallelism = execution (doing multiple things at once)")
	fmt.Println()

	data := []int{16, 4, 10, 14, 7, 9, 3, 2, 8, 1, 15, 11, 12, 6, 5, 13}

	// Sequential sort
	seqData := make([]int, len(data))
	copy(seqData, data)
	seqResult := sequentialMergeSort(seqData)
	fmt.Printf("Sequential sort: %v\n", seqResult)
	fmt.Printf("Sorted correctly: %v\n", isSorted(seqResult))
	fmt.Println()

	// Concurrent sort
	concData := make([]int, len(data))
	copy(concData, data)
	atomic.StoreInt64(&goroutinesSpawned, 0)
	concResult := concurrentMergeSort(concData)
	spawned := atomic.LoadInt64(&goroutinesSpawned)

	fmt.Printf("Concurrent sort: %v\n", concResult)
	fmt.Printf("Sorted correctly: %v\n", isSorted(concResult))
	fmt.Printf("Goroutines spawned: %d\n", spawned)
	fmt.Println()

	fmt.Println("=== Mistake #56: Concurrency Is Not Always Faster ===")
	if spawned > 10 {
		fmt.Printf("WASTEFUL: %d goroutines for %d elements!\n", spawned, len(data))
		fmt.Println("Each goroutine costs ~2KB stack + scheduling overhead.")
		fmt.Println("FIX: Add a threshold to avoid goroutines for small slices.")
	} else {
		fmt.Printf("EFFICIENT: Only %d goroutines for %d elements.\n", spawned, len(data))
		fmt.Println("Goroutines are only used where the work justifies the overhead.")
	}
}
