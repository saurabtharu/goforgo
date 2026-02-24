// test_organization_test.go
// Solution for Mistakes #82, #84, #85
//
// Fixed:
// 1. Table-driven tests replace individual test functions (Mistake #85)
// 2. Subtests with t.Run for clear naming and selective execution (Mistake #84)
// 3. Consolidated duplicate logic into single test functions (Mistake #82)

package main

import "testing"

// FIXED: Single table-driven test replaces TestCategorizeChild, TestCategorizeTeenager,
// TestCategorizeAdult, and TestCategorizeSenior.
func TestCategorize(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected string
	}{
		{"negative age", -1, "invalid"},
		{"zero", 0, "child"},
		{"child", 5, "child"},
		{"child boundary", 12, "child"},
		{"teenager lower", 13, "teenager"},
		{"teenager", 15, "teenager"},
		{"teenager boundary", 17, "teenager"},
		{"adult lower", 18, "adult"},
		{"adult", 30, "adult"},
		{"adult boundary", 64, "adult"},
		{"senior lower", 65, "senior"},
		{"senior", 70, "senior"},
		{"very old", 100, "senior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Categorize(tt.age)
			if result != tt.expected {
				t.Errorf("Categorize(%d) = %q, want %q", tt.age, result, tt.expected)
			}
		})
	}
}

// FIXED: Single table-driven test replaces four separate FizzBuzz test functions.
func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"fizz", 3, "Fizz"},
		{"buzz", 5, "Buzz"},
		{"fizzbuzz", 15, "FizzBuzz"},
		{"number", 7, "7"},
		{"fizz again", 9, "Fizz"},
		{"buzz again", 10, "Buzz"},
		{"fizzbuzz again", 30, "FizzBuzz"},
		{"one", 1, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FizzBuzz(tt.input)
			if result != tt.expected {
				t.Errorf("FizzBuzz(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// FIXED: Single table-driven test replaces four separate ValidateEmail test functions.
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "user@example.com", true},
		{"empty string", "", false},
		{"no at sign", "userexample.com", false},
		{"no domain", "user@", false},
		{"no dot in domain", "user@example", false},
		{"domain too short", "user@ab", false},
		{"at start", "@example.com", false},
		{"valid with subdomain", "user@mail.example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, result, tt.expected)
			}
		})
	}
}
