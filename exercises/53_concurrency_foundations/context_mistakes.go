// context_mistakes.go
// Common context misuse patterns in Go.
// Based on 100 Go Mistakes #60.
//
// Contexts carry deadlines, cancellation signals, and request-scoped values
// across API boundaries. Misusing them leads to leaked goroutines,
// unresponsive services, and hard-to-debug issues.
//
// I AM NOT DONE YET!

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// === Mistake 1: Not propagating context ===
// This handler creates a NEW context.Background() instead of using
// the one passed to it. If the caller cancels, this function won't notice.
//
// FIX: Use the parent ctx parameter instead of context.Background().

func fetchData(ctx context.Context, query string) (string, error) {
	// BUG: Creates a brand new context, ignoring the parent's cancellation!
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simulate slow operation
	// TODO: Check ctx.Done() to respect the caller's cancellation.
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("result for %q", query), nil
}

func handleRequest(ctx context.Context) {
	// The caller may cancel this context, but fetchData ignores it.
	result, err := fetchData(ctx, "users")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	fmt.Printf("  Got: %s\n", result)
}

// === Mistake 2: Not checking context cancellation in loops ===
// This function processes items in a loop but never checks if the
// context has been cancelled. If the caller times out, this keeps running.
//
// FIX: Add a select{} inside the loop that checks ctx.Done() on each iteration.

func processItems(ctx context.Context, items []string) []string {
	var results []string
	for _, item := range items {
		// TODO: Check ctx.Done() before processing each item.
		// Use a select with ctx.Done() and a default case.

		// Simulate work per item
		time.Sleep(50 * time.Millisecond)
		results = append(results, fmt.Sprintf("processed:%s", item))
	}
	return results
}

// === Mistake 3: Deadline vs Timeout confusion ===
// WithTimeout(ctx, 5s) sets a deadline 5s from NOW.
// WithDeadline(ctx, t) sets an absolute deadline at time t.
// Nesting timeouts: the shorter one always wins.
//
// FIX: The inner function uses a 10-second timeout, but the outer context
// has a 200ms deadline. The inner timeout is pointless.
// Change the inner function to NOT create its own timeout - just use the parent ctx.

func outerHandler() {
	// Outer: 200ms deadline
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := innerWork(ctx)
	fmt.Printf("  Result: %s\n", result)
}

func innerWork(ctx context.Context) string {
	// BUG: This 10-second timeout is pointless - the parent's 200ms deadline
	// is shorter and will fire first. This just wastes a goroutine on the timer.
	innerCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-time.After(150 * time.Millisecond):
		return "work completed"
	case <-innerCtx.Done():
		return fmt.Sprintf("cancelled: %v", innerCtx.Err())
	}
}

// === Mistake 4: Using context values for the wrong things ===
// Context values should only carry request-scoped data that crosses
// API boundaries (trace IDs, auth tokens, request IDs).
// Do NOT use context to pass function parameters or configuration.
//
// FIX: Move the "format" and "limit" values from context to function parameters.
// Keep only the requestID in context (it's genuinely request-scoped).

type contextKey string

const (
	keyRequestID contextKey = "requestID"
	keyFormat    contextKey = "format"
	keyLimit     contextKey = "limit"
)

// BUG: format and limit are function parameters disguised as context values.
// They should be regular function arguments.
func formatResults(ctx context.Context, data []string) string {
	requestID, _ := ctx.Value(keyRequestID).(string)
	format, _ := ctx.Value(keyFormat).(string)
	limit, _ := ctx.Value(keyLimit).(int)

	if limit <= 0 || limit > len(data) {
		limit = len(data)
	}

	result := fmt.Sprintf("[req:%s] ", requestID)
	for i := 0; i < limit; i++ {
		if format == "upper" {
			result += fmt.Sprintf("%d:%s ", i, data[i])
		} else {
			result += fmt.Sprintf("%d:%s ", i, data[i])
		}
	}
	return result
}

func main() {
	fmt.Println("=== Mistake #60: Context Misuse Patterns ===")
	fmt.Println()

	// Mistake 1: Not propagating context
	fmt.Println("--- Mistake 1: Not Propagating Context ---")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	handleRequest(ctx)
	// With the fix, fetchData should respect the 50ms cancellation.
	// Without the fix, it ignores it and always completes.
	fmt.Println()

	// Mistake 2: Not checking cancellation in loops
	fmt.Println("--- Mistake 2: Not Checking Cancellation ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel2()
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	results := processItems(ctx2, items)
	fmt.Printf("  Processed %d/%d items\n", len(results), len(items))
	if len(results) < len(items) {
		fmt.Println("  OK: Stopped early due to context cancellation")
	} else {
		fmt.Println("  BUG: Processed all items ignoring cancellation!")
	}
	fmt.Println()

	// Mistake 3: Deadline vs Timeout
	fmt.Println("--- Mistake 3: Deadline vs Timeout ---")
	outerHandler()
	fmt.Println()

	// Mistake 4: Context values misuse
	fmt.Println("--- Mistake 4: Context Values Misuse ---")
	ctx4 := context.WithValue(context.Background(), keyRequestID, "req-123")
	ctx4 = context.WithValue(ctx4, keyFormat, "upper")
	ctx4 = context.WithValue(ctx4, keyLimit, 3)
	data := []string{"alpha", "beta", "gamma", "delta"}
	formatted := formatResults(ctx4, data)
	fmt.Printf("  %s\n", formatted)
	fmt.Println("  TODO: Move format and limit to function parameters")

	// Summary
	fmt.Println()
	fmt.Println("=== Context Best Practices ===")
	fmt.Println("1. Always propagate the parent context")
	fmt.Println("2. Check ctx.Done() in long-running loops")
	fmt.Println("3. Understand that shorter deadlines always win")
	fmt.Println("4. Only use context values for request-scoped data")

	// Keep references to avoid unused import errors
	var _ sync.Mutex
}
