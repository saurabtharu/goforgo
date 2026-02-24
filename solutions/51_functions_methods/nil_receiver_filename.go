package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// --- Mistake #45: Returning a nil receiver wrapped in an interface ---

// ValidationError is a custom error type.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s - %s", e.Field, e.Message)
}

// validate checks if a name is valid.
// Fixed: Return nil directly instead of a typed nil pointer.
// When the interface itself is nil, `if err != nil` works correctly.
func validate(name string) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	// Fixed: return nil directly, not a (*ValidationError)(nil) typed pointer.
	// A nil interface value has both its type and value as nil.
	return nil
}

// --- Mistake #46: Using filename instead of io.Reader ---

// countWordsFromReader accepts an io.Reader, making it flexible and testable.
// You can pass a file, a string reader, an HTTP response body, or any reader.
func countWordsFromReader(r io.Reader) (int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, fmt.Errorf("reading: %w", err)
	}
	words := strings.Fields(string(data))
	return len(words), nil
}

// countWords is now a thin wrapper that opens a file and delegates to countWordsFromReader.
// The core logic is in countWordsFromReader where it can be easily tested.
func countWords(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()
	return countWordsFromReader(f)
}

func main() {
	// Fixed: validate returns a true nil interface for valid names
	err := validate("Alice")
	if err != nil {
		fmt.Printf("Unexpected error for valid name: %v\n", err)
	} else {
		fmt.Println("Alice: valid")
	}

	err = validate("")
	if err != nil {
		fmt.Printf("Empty name: %v\n", err)
	}

	// Fixed: Using io.Reader makes this testable without real files
	reader := strings.NewReader("the quick brown fox jumps over the lazy dog")
	count, err := countWordsFromReader(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Word count: %d\n", count)
	}
}
