// context_propagation.go
// Learn why using a request context in background goroutines is dangerous
//
// Mistake #61 from "100 Go Mistakes": Propagating an inappropriate context
// When an HTTP handler spawns a background goroutine using the request context,
// the context gets canceled when the handler returns, killing the background work.

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
	reqCtx, reqCancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// BUG: This background goroutine uses reqCtx directly.
		// When the handler "returns" and calls reqCancel(), this context
		// is canceled and the background work is aborted.
		//
		// TODO: Fix this by creating a detached context (context.Background()
		// with its own timeout) instead of using reqCtx for the background work.
		// The background task needs 200ms to complete, so give it a 500ms timeout.
		err := doBackgroundWork(reqCtx, "audit-log-write")
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
