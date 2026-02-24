package main

import (
	"errors"
	"fmt"
	"strings"
)

// 100 Go Mistakes #52, #53, #54: Error Handling Patterns
//
// Three bugs to fix:
// 1. (#52) processOrder logs the error AND returns it — the caller also
//    logs it, producing duplicate messages.
// 2. (#53) processItems ignores the error from shipItem — failed shipments
//    go unnoticed.
// 3. (#54) saveReport uses defer w.Close() which silently drops the
//    Close() error — data may not be flushed and the caller never knows.
//
// FIX:
// 1. Handle errors once: either log or return, not both
// 2. Always handle returned errors
// 3. Use named return values to capture defer errors

// --- Simulated logger ---

var logMessages []string

func logError(msg string) {
	logMessages = append(logMessages, msg)
}

// --- Bug 1: Handling error twice ---

func validateOrder(id string) error {
	if id == "" {
		return errors.New("order ID is empty")
	}
	return nil
}

func processOrder(orderID string) error {
	err := validateOrder(orderID)
	if err != nil {
		// BUG: Logs the error AND returns it.
		// The caller (handleOrder) will also log it → duplicate messages.
		logError(fmt.Sprintf("processOrder failed: %v", err))
		return fmt.Errorf("processOrder: %w", err)
	}
	return nil
}

func handleOrder(orderID string) {
	err := processOrder(orderID)
	if err != nil {
		// This is the second time the error gets logged!
		logError(fmt.Sprintf("handleOrder failed: %v", err))
	}
}

// --- Bug 2: Not handling errors ---

func shipItem(item string) error {
	if strings.HasPrefix(item, "fragile") {
		return fmt.Errorf("special handling required for %q", item)
	}
	fmt.Printf("Shipped: %s\n", item)
	return nil
}

func processItems(items []string) {
	for _, item := range items {
		// BUG: Error return value is silently ignored!
		// If shipping fails, we'll never know.
		shipItem(item)
	}
	fmt.Println("All items processed")
}

// --- Bug 3: Not handling defer errors ---

type DataWriter struct {
	name   string
	buffer []string
}

func NewDataWriter(name string) *DataWriter {
	return &DataWriter{name: name}
}

func (w *DataWriter) Write(data string) {
	w.buffer = append(w.buffer, data)
}

func (w *DataWriter) Close() error {
	if len(w.buffer) == 0 {
		return fmt.Errorf("nothing to flush")
	}
	fmt.Printf("Writer %q: flushed %d items\n", w.name, len(w.buffer))
	return nil
}

func saveReport(name string, lines []string) error {
	w := NewDataWriter(name)
	// BUG: Close() error is silently dropped by defer!
	// If flush fails, the caller will never know data wasn't saved.
	defer w.Close()

	for _, line := range lines {
		w.Write(line)
	}
	return nil
}

func main() {
	fmt.Println("=== Error Handling Patterns ===")

	// Bug 1: Double handling
	fmt.Println("--- Bug 1: Handling Error Twice ---")
	logMessages = nil
	handleOrder("")
	fmt.Printf("Log entries: %d (should be 1, not 2!)\n", len(logMessages))
	for i, msg := range logMessages {
		fmt.Printf("  [%d] %s\n", i+1, msg)
	}

	fmt.Println()

	// Bug 2: Ignoring errors
	fmt.Println("--- Bug 2: Ignoring Errors ---")
	processItems([]string{"book", "fragile-vase", "laptop"})

	fmt.Println()

	// Bug 3: Defer error dropped
	fmt.Println("--- Bug 3: Lost Defer Error ---")
	err := saveReport("empty-report", nil)
	fmt.Printf("saveReport error: %v\n", err)

	err = saveReport("monthly", []string{"revenue: $1000", "expenses: $800"})
	fmt.Printf("saveReport error: %v\n", err)
}
