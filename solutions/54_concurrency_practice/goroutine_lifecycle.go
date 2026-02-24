// goroutine_lifecycle.go - SOLUTION
// Learn to manage goroutine lifetimes and avoid common goroutine leaks
//
// Mistakes #62-63 from "100 Go Mistakes":
// #62: Goroutines that are started but never stopped (goroutine leaks)
// #63: Closure variable capture in goroutines

package main

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"sync"
)

// === Part 1: Goroutine Leak (Mistake #62) ===

// stoppableProducer sends values on a channel but respects context cancellation.
// When the context is canceled, the goroutine exits cleanly.
func stoppableProducer(ctx context.Context, ch chan<- int) {
	i := 0
	for {
		select {
		case ch <- i:
			i++
		case <-ctx.Done():
			return
		}
	}
}

// consumeN reads exactly n values from the channel and returns them.
func consumeN(ch <-chan int, n int) []int {
	results := make([]int, 0, n)
	for i := 0; i < n; i++ {
		results = append(results, <-ch)
	}
	return results
}

// === Part 2: Closure Variable Capture (Mistake #63) ===

// FIX: Pass the value as a function parameter so each goroutine gets its own copy
func processItems(items []string) []string {
	ch := make(chan string, len(items))
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		// FIX: Pass item as parameter - each goroutine gets its own copy
		go func(val string) {
			defer wg.Done()
			ch <- fmt.Sprintf("processed: %s", val)
		}(item)
	}

	wg.Wait()
	close(ch)

	var results []string
	for result := range ch {
		results = append(results, result)
	}
	sort.Strings(results)
	return results
}

func main() {
	fmt.Println("=== Goroutine Leak (Mistake #62) ===")

	goroutinesBefore := runtime.NumGoroutine()

	// FIX: Create a cancelable context to control the producer lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int, 1)

	// FIX: Pass context so the producer can be stopped
	go stoppableProducer(ctx, ch)

	values := consumeN(ch, 5)
	fmt.Printf("consumed: %v\n", values)

	// FIX: Cancel context to stop the producer goroutine cleanly
	cancel()

	runtime.Gosched() // Let goroutines settle
	goroutinesAfter := runtime.NumGoroutine()
	leaked := goroutinesAfter - goroutinesBefore
	if leaked > 0 {
		fmt.Printf("FAIL: %d goroutine(s) leaked! Producer is still running.\n", leaked)
	} else {
		fmt.Println("PASS: No goroutine leaks - producer stopped cleanly")
	}

	fmt.Println()
	fmt.Println("=== Closure Variable Capture (Mistake #63) ===")

	items := []string{"alpha", "beta", "gamma", "delta"}
	results := processItems(items)

	fmt.Println("expected: each item processed exactly once")
	for _, r := range results {
		fmt.Println(r)
	}

	// Verify all items were processed (not duplicates of the last one)
	expected := []string{
		"processed: alpha",
		"processed: beta",
		"processed: delta",
		"processed: gamma",
	}
	allCorrect := len(results) == len(expected)
	if allCorrect {
		for i, r := range results {
			if r != expected[i] {
				allCorrect = false
				break
			}
		}
	}
	if allCorrect {
		fmt.Println("PASS: All items processed correctly")
	} else {
		fmt.Println("FAIL: Closure variable was captured incorrectly")
	}
}
