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
// BUG: This returns a *ValidationError typed nil, not a nil interface.
// When the name is valid, it returns (*ValidationError)(nil), which when
// assigned to an error interface, makes the interface non-nil!
// FIX: Return nil directly instead of the typed nil pointer.
func validate(name string) error {
	var err *ValidationError
	if name == "" {
		err = &ValidationError{Field: "name", Message: "cannot be empty"}
	}
	// BUG: even when name is valid, err is (*ValidationError)(nil),
	// and returning it wraps it in a non-nil error interface
	return err
}

// --- Mistake #46: Using filename instead of io.Reader ---

// countWords takes a filename and counts words in the file.
// FIX: This function is hard to test because it requires a real file on disk.
// Refactor it to accept an io.Reader instead, making it flexible and testable.
// Then create a wrapper function that opens the file and calls the reader version.
func countWords(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("reading file: %w", err)
	}
	words := strings.Fields(string(data))
	return len(words), nil
}

func main() {
	// Demonstrate the nil interface gotcha
	err := validate("Alice")
	if err != nil {
		// BUG: This branch runs even though "Alice" is a valid name!
		// That's because validate returns (*ValidationError)(nil) wrapped
		// in the error interface, making the interface non-nil.
		fmt.Printf("Unexpected error for valid name: %v\n", err)
	} else {
		fmt.Println("Alice: valid")
	}

	err = validate("")
	if err != nil {
		fmt.Printf("Empty name: %v\n", err)
	}

	// Demonstrate io.Reader usage
	// FIX: After refactoring countWords to accept io.Reader,
	// use strings.NewReader here instead of writing to a temp file.
	reader := strings.NewReader("the quick brown fox jumps over the lazy dog")
	count, err := countWordsFromReader(reader)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Word count: %d\n", count)
	}
}

// TODO: Implement countWordsFromReader that accepts an io.Reader.
// This is the testable version that countWords should delegate to.
func countWordsFromReader(r io.Reader) (int, error) {
	// FIX: Implement this function. Read all bytes from r,
	// split into words, and return the count.
	_ = r
	return 0, fmt.Errorf("not implemented")
}
