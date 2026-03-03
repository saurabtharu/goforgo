package main

import (
	"errors"
	"fmt"
)

// 100 Go Mistakes #50, #51: Error Comparison (Solution)
//
// Fixed: Use errors.As for type checks and errors.Is for sentinel checks.
// Both functions walk the entire wrapping chain, finding matches that
// type assertions and == would miss.

// --- Custom error types ---

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: field %q %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %d not found", e.Resource, e.ID)
}

// --- Sentinel errors ---

var (
	ErrPermission = errors.New("permission denied")
	ErrRateLimit  = errors.New("rate limit exceeded")
)

// --- Functions that return WRAPPED errors ---

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return fmt.Errorf("validateAge: %w", &ValidationError{
			Field: "age", Message: "must be between 0 and 150",
		})
	}
	return nil
}

func findUser(id int) error {
	if id == 99 {
		return fmt.Errorf("findUser: %w", &NotFoundError{
			Resource: "user", ID: id,
		})
	}
	return nil
}

func checkAccess(role string) error {
	if role != "admin" {
		return fmt.Errorf("checkAccess role=%q: %w", role, ErrPermission)
	}
	return nil
}

func callAPI() error {
	return fmt.Errorf("callAPI: %w", ErrRateLimit)
}

func main() {
	fmt.Println("=== Error Comparison ===")

	// Fix #50: Use errors.As to extract error types through wrapping chain
	fmt.Println("--- Fix #50: errors.As Unwraps the Chain ---")

	err := validateAge(200)
	// Fixed: errors.As walks the chain and finds *ValidationError inside.
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("Validation error: field=%q msg=%q\n", ve.Field, ve.Message)
	}

	err = findUser(99)
	// Fixed: errors.As finds *NotFoundError even though it's wrapped.
	var nfe *NotFoundError
	if errors.As(err, &nfe) {
		fmt.Printf("Not found: %s id=%d\n", nfe.Resource, nfe.ID)
	}

	fmt.Println()

	// Fix #51: Use errors.Is to check sentinel errors through wrapping chain
	fmt.Println("--- Fix #51: errors.Is Unwraps the Chain ---")

	err = checkAccess("viewer")
	// Fixed: errors.Is walks the chain and finds ErrPermission inside.
	if errors.Is(err, ErrPermission) {
		fmt.Println("Permission denied - access rejected")
	}

	err = callAPI()
	// Fixed: errors.Is finds ErrRateLimit even through wrapping.
	if errors.Is(err, ErrRateLimit) {
		fmt.Println("Rate limited - will retry")
	}
}
