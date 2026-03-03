package main

import (
	"errors"
	"fmt"
)

// 100 Go Mistakes #50, #51: Error Comparison
//
// Two common mistakes when checking errors:
// - #50: Using type assertion (err.(*Type)) on wrapped errors — fails
//        because the outer error is *fmt.wrapError, not your custom type.
// - #51: Using == comparison (err == sentinel) on wrapped errors — fails
//        because the outer error is a different object than the sentinel.
//
// FIX: Replace type assertions with errors.As() and == with errors.Is().
// Both functions walk the entire wrapping chain to find a match.

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

	// --- Mistake #50: Type assertions fail on wrapped errors ---
	fmt.Println("--- Mistake #50: Type Assertion on Wrapped Errors ---")

	err := validateAge(200)
	// BUG: Type assertion fails — err is *fmt.wrapError, not *ValidationError.
	// The ValidationError is inside the wrapping chain, but .(type) only
	// checks the outermost error.
	if ve, ok := err.(*ValidationError); ok {
		fmt.Printf("Validation error: field=%q msg=%q\n", ve.Field, ve.Message)
	} else {
		fmt.Println("FAIL: type assertion missed ValidationError")
	}

	err = findUser(99)
	// BUG: Same problem — type assertion can't see through wrapping.
	if nfe, ok := err.(*NotFoundError); ok {
		fmt.Printf("Not found: %s id=%d\n", nfe.Resource, nfe.ID)
	} else {
		fmt.Println("FAIL: type assertion missed NotFoundError")
	}

	fmt.Println()

	// --- Mistake #51: == comparison fails on wrapped sentinel errors ---
	fmt.Println("--- Mistake #51: == Comparison on Wrapped Errors ---")

	err = checkAccess("viewer")
	// BUG: == can't see through wrapping. The outer error is a *fmt.wrapError,
	// not ErrPermission itself.
	if err == ErrPermission {
		fmt.Println("Permission denied - access rejected")
	} else {
		fmt.Println("FAIL: == missed ErrPermission")
	}

	err = callAPI()
	// BUG: Same problem with ==.
	if err == ErrRateLimit {
		fmt.Println("Rate limited - will retry")
	} else {
		fmt.Println("FAIL: == missed ErrRateLimit")
	}
}
