// test_organization.go
// Mistakes #82, #84, #85: Test categorization, execution flags, and table-driven tests
//
// This file contains the functions to be tested.
// Your task is to fix the tests in test_organization_test.go.

package main

import (
	"fmt"
	"strings"
)

// Categorize returns an age group label for a given age.
func Categorize(age int) string {
	switch {
	case age < 0:
		return "invalid"
	case age < 13:
		return "child"
	case age < 18:
		return "teenager"
	case age < 65:
		return "adult"
	default:
		return "senior"
	}
}

// FizzBuzz returns "Fizz" for multiples of 3, "Buzz" for multiples of 5,
// "FizzBuzz" for multiples of both, or the number as a string.
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprintf("%d", n)
	}
}

// ValidateEmail performs a basic email validation check.
func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex < 1 {
		return false
	}
	domain := email[atIndex+1:]
	if len(domain) < 3 {
		return false
	}
	if !strings.Contains(domain, ".") {
		return false
	}
	return true
}

func main() {
	fmt.Println("Run 'go test -v' to execute the tests.")
}
