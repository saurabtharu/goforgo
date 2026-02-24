package main

import (
	"errors"
	"fmt"
	"strings"
)

// 100 Go Mistakes #52, #53, #54: Error Handling Patterns (Solution)
//
// Fixed all three anti-patterns:
// 1. Handle errors once (return, don't also log)
// 2. Always handle returned errors
// 3. Use named returns to capture defer errors

// --- Simulated logger ---

var logMessages []string

func logError(msg string) {
	logMessages = append(logMessages, msg)
}

// --- Fix 1: Handle error once ---

func validateOrder(id string) error {
	if id == "" {
		return errors.New("order ID is empty")
	}
	return nil
}

func processOrder(orderID string) error {
	err := validateOrder(orderID)
	if err != nil {
		// Fixed: Just return with context, don't also log.
		// The caller is responsible for deciding how to handle it.
		return fmt.Errorf("processOrder: %w", err)
	}
	return nil
}

func handleOrder(orderID string) {
	err := processOrder(orderID)
	if err != nil {
		// This is the ONE place where the error is handled (logged).
		logError(fmt.Sprintf("handleOrder failed: %v", err))
	}
}

// --- Fix 2: Handle errors from shipItem ---

func shipItem(item string) error {
	if strings.HasPrefix(item, "fragile") {
		return fmt.Errorf("special handling required for %q", item)
	}
	fmt.Printf("Shipped: %s\n", item)
	return nil
}

func processItems(items []string) {
	shipped, failed := 0, 0
	for _, item := range items {
		// Fixed: Capture and handle the error.
		err := shipItem(item)
		if err != nil {
			fmt.Printf("Error shipping %q: %v\n", item, err)
			failed++
		} else {
			shipped++
		}
	}
	fmt.Printf("Processed %d items: %d shipped, %d failed\n", len(items), shipped, failed)
}

// --- Fix 3: Use named return to capture defer error ---

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

// Fixed: Named return value (err) lets the defer closure update the
// return value if Close() fails.
func saveReport(name string, lines []string) (err error) {
	w := NewDataWriter(name)
	defer func() {
		closeErr := w.Close()
		if closeErr != nil {
			if err != nil {
				// Both the body and Close() failed — combine them.
				err = fmt.Errorf("%w; close %q: %v", err, name, closeErr)
			} else {
				// Body succeeded but Close() failed — report it.
				err = fmt.Errorf("close %q: %v", name, closeErr)
			}
		}
	}()

	for _, line := range lines {
		w.Write(line)
	}
	return nil
}

func main() {
	fmt.Println("=== Error Handling Patterns ===")

	// Fix 1: Single handling point
	fmt.Println("--- Fix 1: Handle Error Once ---")
	logMessages = nil
	handleOrder("")
	fmt.Printf("Log entries: %d\n", len(logMessages))
	for i, msg := range logMessages {
		fmt.Printf("  [%d] %s\n", i+1, msg)
	}

	fmt.Println()

	// Fix 2: All errors handled
	fmt.Println("--- Fix 2: Handle All Errors ---")
	processItems([]string{"book", "fragile-vase", "laptop"})

	fmt.Println()

	// Fix 3: Defer error captured
	fmt.Println("--- Fix 3: Handle Defer Error ---")
	err := saveReport("empty-report", nil)
	fmt.Printf("saveReport error: %v\n", err)

	err = saveReport("monthly", []string{"revenue: $1000", "expenses: $800"})
	fmt.Printf("saveReport error: %v\n", err)
}
