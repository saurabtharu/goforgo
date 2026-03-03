// context_propagation.go - SOLUTION
// Learn why using a request context in background goroutines is dangerous
//
// Mistake #61 from "100 Go Mistakes": Propagating an inappropriate context

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// simulateHTTPHandler pretends to be an HTTP handler that accepts a request,
// kicks off background work, and then "returns" (canceling its context).
func simulateHTTPHandler(results chan<- string) {
	// This context simulates the request-scoped context that net/http provides.
	// It gets canceled when the handler returns.
	_, reqCancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// FIX: Create a detached context for background work.
		// The background task should NOT be tied to the request lifecycle.
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer bgCancel()

		err := doBackgroundWork(bgCtx, "audit-log-write")
		if err != nil {
			results <- fmt.Sprintf("background work failed: %v", err)
		} else {
			results <- "background work succeeded"
		}
	}()

	// Simulate the handler doing its main work quickly
	fmt.Println("handler: processing request...")
	time.Sleep(50 * time.Millisecond)
	fmt.Println("handler: sending response to client")

	// Handler returns, canceling the request context.
	// In real net/http, this happens automatically.
	reqCancel()
	fmt.Println("handler: context canceled (handler returned)")

	wg.Wait()
}

// doBackgroundWork simulates work that takes 200ms.
// It respects context cancellation.
func doBackgroundWork(ctx context.Context, taskName string) error {
	fmt.Printf("background: starting %s\n", taskName)

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Printf("background: %s completed\n", taskName)
		return nil
	case <-ctx.Done():
		fmt.Printf("background: %s canceled: %v\n", taskName, ctx.Err())
		return ctx.Err()
	}
}

func main() {
	fmt.Println("=== Context Propagation (Mistake #61) ===")
	fmt.Println("Demonstrating the danger of passing request contexts to background goroutines")
	fmt.Println()

	results := make(chan string, 1)
	simulateHTTPHandler(results)
	result := <-results

	fmt.Println()
	if result == "background work succeeded" {
		fmt.Println("PASS: Background work completed despite handler returning")
	} else {
		fmt.Printf("FAIL: %s\n", result)
		fmt.Println("Hint: The background goroutine should NOT use the request context")
	}
}
