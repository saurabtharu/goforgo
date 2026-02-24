// test_organization_test.go
// Mistakes #82, #84, #85: Fix these tests!
//
// Problems to fix:
// 1. Individual test functions that should be table-driven tests (Mistake #85)
// 2. No subtests with t.Run for clear test naming (Mistake #84)
// 3. Duplicate test logic that should be consolidated (Mistake #82)
//
// Refactor all tests below into three table-driven test functions:
//   - TestCategorize (with subtests via t.Run)
//   - TestFizzBuzz (with subtests via t.Run)
//   - TestValidateEmail (with subtests via t.Run)
//
// Then delete the TestRefactorToTableDriven function at the bottom.

package main

import "testing"

// BUG #1: These are separate test functions testing the same function
// with different inputs. They should be a single table-driven test.

func TestCategorizeChild(t *testing.T) {
	result := Categorize(5)
	if result != "child" {
		t.Errorf("Categorize(5) = %q, want %q", result, "child")
	}
}

func TestCategorizeTeenager(t *testing.T) {
	result := Categorize(15)
	if result != "teenager" {
		t.Errorf("Categorize(15) = %q, want %q", result, "teenager")
	}
}

func TestCategorizeAdult(t *testing.T) {
	result := Categorize(30)
	if result != "adult" {
		t.Errorf("Categorize(30) = %q, want %q", result, "adult")
	}
}

func TestCategorizeSenior(t *testing.T) {
	result := Categorize(70)
	if result != "senior" {
		t.Errorf("Categorize(70) = %q, want %q", result, "senior")
	}
}

// BUG #2: FizzBuzz tests also have duplicated logic and no subtests.
// Each case is copy-pasted with only the values changing.

func TestFizzBuzzFizz(t *testing.T) {
	result := FizzBuzz(3)
	if result != "Fizz" {
		t.Errorf("FizzBuzz(3) = %q, want %q", result, "Fizz")
	}
}

func TestFizzBuzzBuzz(t *testing.T) {
	result := FizzBuzz(5)
	if result != "Buzz" {
		t.Errorf("FizzBuzz(5) = %q, want %q", result, "Buzz")
	}
}

func TestFizzBuzzFizzBuzz(t *testing.T) {
	result := FizzBuzz(15)
	if result != "FizzBuzz" {
		t.Errorf("FizzBuzz(15) = %q, want %q", result, "FizzBuzz")
	}
}

func TestFizzBuzzNumber(t *testing.T) {
	result := FizzBuzz(7)
	if result != "7" {
		t.Errorf("FizzBuzz(7) = %q, want %q", result, "7")
	}
}

// BUG #3: ValidateEmail tests should also be table-driven.
// The repetitive assertion pattern is a clear signal.

func TestValidateEmailValid(t *testing.T) {
	result := ValidateEmail("user@example.com")
	if !result {
		t.Error("ValidateEmail(\"user@example.com\") should be true")
	}
}

func TestValidateEmailEmpty(t *testing.T) {
	result := ValidateEmail("")
	if result {
		t.Error("ValidateEmail(\"\") should be false")
	}
}

func TestValidateEmailNoAt(t *testing.T) {
	result := ValidateEmail("userexample.com")
	if result {
		t.Error("ValidateEmail(\"userexample.com\") should be false")
	}
}

func TestValidateEmailNoDomain(t *testing.T) {
	result := ValidateEmail("user@")
	if result {
		t.Error("ValidateEmail(\"user@\") should be false")
	}
}

// This test forces the exercise to fail until you refactor.
// Delete this function after converting all tests above to table-driven tests.
func TestRefactorToTableDriven(t *testing.T) {
	t.Fatal("EXERCISE: Refactor all tests above into 3 table-driven test functions " +
		"(TestCategorize, TestFizzBuzz, TestValidateEmail) using t.Run subtests, " +
		"then delete this function.")
}
