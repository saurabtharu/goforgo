// goroutine_lifecycle.go
// Learn to manage goroutine lifetimes and avoid common goroutine leaks
//
// Mistakes #62-63 from "100 Go Mistakes":
// #62: Goroutines that are started but never stopped (goroutine leaks)
// #63: Closure variable capture in goroutines

package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
)

// NOTE: You'll need to import "context" to fix Part 1

// === Part 1: Goroutine Leak (Mistake #62) ===

// leakyProducer sends values on a channel forever.
// BUG: If the consumer stops reading, this goroutine blocks forever on send,
// leaking the goroutine. There's no way to signal it to stop.
//
// TODO: Add a context.Context parameter and use select to check for
// cancellation when sending on the channel, so the goroutine can be stopped.
func leakyProducer(ch chan<- int) {
	i := 0
	for {
		ch <- i // Blocks forever if nobody reads
		i++
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

// processItems launches a goroutine for each item but captures a shared
// variable incorrectly.
//
// BUG: The variable 'current' is declared once outside the loop. All
// goroutines capture the same variable and see whatever value it holds
// when they eventually execute (usually the last value).
//
// TODO: Fix by passing 'current' as a goroutine function parameter,
// giving each goroutine its own copy.
func processItems(items []string) []string {
	ch := make(chan string, len(items))
	var wg sync.WaitGroup

	current := "" // BUG: single variable shared across all goroutine closures
	for _, item := range items {
		current = item // all goroutines will see the final value of current
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- fmt.Sprintf("processed: %s", current)
		}()
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

	// TODO: Create a context you can cancel to stop the producer
	ch := make(chan int, 1)

	// TODO: Pass context to leakyProducer so it can be stopped
	go leakyProducer(ch)

	values := consumeN(ch, 5)
	fmt.Printf("consumed: %v\n", values)

	// TODO: Cancel the context here to stop the producer goroutine
	// Without cancellation, the producer goroutine leaks!

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
