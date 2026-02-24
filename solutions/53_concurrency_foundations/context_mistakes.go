// context_mistakes.go - SOLUTION
// Common context misuse patterns in Go.
// Based on 100 Go Mistakes #60.

package main

import (
	"context"
	"fmt"
	"time"
)

// === Mistake 1: FIXED - propagate the parent context ===

func fetchData(ctx context.Context, query string) (string, error) {
	// Use the parent context so cancellation propagates correctly.
	select {
	case <-time.After(100 * time.Millisecond):
		return fmt.Sprintf("result for %q", query), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func handleRequest(ctx context.Context) {
	result, err := fetchData(ctx, "users")
	if err != nil {
		fmt.Printf("  Cancelled: %v\n", err)
		return
	}
	fmt.Printf("  Got: %s\n", result)
}

// === Mistake 2: FIXED - check cancellation in the loop ===

func processItems(ctx context.Context, items []string) []string {
	var results []string
	for _, item := range items {
		// Check context before doing work on each item.
		select {
		case <-ctx.Done():
			return results
		default:
		}

		time.Sleep(50 * time.Millisecond)
		results = append(results, fmt.Sprintf("processed:%s", item))
	}
	return results
}

// === Mistake 3: FIXED - don't create redundant inner timeout ===

func outerHandler() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := innerWork(ctx)
	fmt.Printf("  Result: %s\n", result)
}

func innerWork(ctx context.Context) string {
	// Just use the parent context directly - no redundant inner timeout.
	select {
	case <-time.After(150 * time.Millisecond):
		return "work completed"
	case <-ctx.Done():
		return fmt.Sprintf("cancelled: %v", ctx.Err())
	}
}

// === Mistake 4: FIXED - use function parameters, not context values ===

type contextKey string

const keyRequestID contextKey = "requestID"

// format and limit are now proper function parameters.
// Only requestID stays in context (genuinely request-scoped).
func formatResults(ctx context.Context, data []string, format string, limit int) string {
	requestID, _ := ctx.Value(keyRequestID).(string)

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

	// Mistake 1: Context now propagated correctly
	fmt.Println("--- Mistake 1: Not Propagating Context ---")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	handleRequest(ctx)
	fmt.Println()

	// Mistake 2: Cancellation now checked in loop
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

	// Mistake 3: No redundant inner timeout
	fmt.Println("--- Mistake 3: Deadline vs Timeout ---")
	outerHandler()
	fmt.Println()

	// Mistake 4: format and limit are now function parameters
	fmt.Println("--- Mistake 4: Context Values Misuse ---")
	ctx4 := context.WithValue(context.Background(), keyRequestID, "req-123")
	data := []string{"alpha", "beta", "gamma", "delta"}
	formatted := formatResults(ctx4, data, "upper", 3)
	fmt.Printf("  %s\n", formatted)
	fmt.Println("  OK: format and limit are function parameters now")

	fmt.Println()
	fmt.Println("=== Context Best Practices ===")
	fmt.Println("1. Always propagate the parent context")
	fmt.Println("2. Check ctx.Done() in long-running loops")
	fmt.Println("3. Understand that shorter deadlines always win")
	fmt.Println("4. Only use context values for request-scoped data")
}
