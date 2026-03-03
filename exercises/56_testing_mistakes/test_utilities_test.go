// test_utilities_test.go
// Mistakes #88, #90, #91: Fix these tests!
//
// Problems to fix:
// 1. HTTP tests don't use httptest (Mistake #88)
//    Use httptest.NewRecorder and http.NewRequest to properly test handlers
// 2. Helper functions don't call t.Helper() (Mistake #90)
//    Error messages point to the helper instead of the calling test
// 3. Tests use defer instead of t.Cleanup for resource cleanup (Mistake #91)
//    t.Cleanup is tied to the test lifecycle and runs even on panic

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// BUG #1 (Mistake #90): This helper function doesn't call t.Helper().
// When it fails, the error points to this function instead of the
// calling test. Add t.Helper() as the first line.
func assertStatusCode(t *testing.T, got, want int) {
	// TODO: Add t.Helper() here so error messages point to the caller
	if got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}
}

// BUG #2 (Mistake #90): Same issue - no t.Helper() call.
func assertContains(t *testing.T, body, substr string) {
	// TODO: Add t.Helper() here
	for i := 0; i <= len(body)-len(substr); i++ {
		if body[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("body %q does not contain %q", body, substr)
}

// BUG #3 (Mistake #88): This test doesn't use httptest properly.
// It creates a recorder and request but has the wrong expected value.
// Fix the assertion to match the actual JSON output from the handler.
func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	HealthHandler(rec, req)

	assertStatusCode(t, rec.Code, http.StatusOK)

	body, _ := io.ReadAll(rec.Body)
	// BUG: Comparing against wrong expected value - json.Encoder adds a newline
	expected := `{"status":"ok"}`
	if string(body) != expected {
		t.Errorf("body = %q, want %q", string(body), expected)
	}
}

// BUG #4 (Mistake #88): UserHandler test is incomplete.
// It only tests the success case. Add subtests for missing id
// and invalid id parameters.
func TestUserHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user?id=1", nil)

	UserHandler(rec, req)

	assertStatusCode(t, rec.Code, http.StatusOK)

	var user User
	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("user.ID = %d, want 1", user.ID)
	}

	// BUG: Missing tests for error cases (missing id, invalid id)
	// TODO: Add t.Run subtests for "missing id" and "invalid id"
}

// BUG #5 (Mistake #88): GreetHandler test is missing entirely.
// TODO: Write a proper test using httptest for both default and custom names.
func TestGreetHandler(t *testing.T) {
	// TODO: Implement this test using httptest.NewRecorder and httptest.NewRequest
	// Test both: GET /greet (should return "Hello, World!")
	// and: GET /greet?name=Gopher (should return "Hello, Gopher!")
	t.Fatal("TODO: Implement TestGreetHandler using httptest")
}

// BUG #6 (Mistake #91): This test creates a temp file but uses defer
// for cleanup instead of t.Cleanup.
func TestWriteAndReadFile(t *testing.T) {
	path, err := WriteTempFile("hello world")
	if err != nil {
		t.Fatalf("WriteTempFile failed: %v", err)
	}
	// BUG: Using defer for cleanup instead of t.Cleanup
	defer os.Remove(path)

	content, err := ReadFileContent(path)
	if err != nil {
		t.Fatalf("ReadFileContent failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("content = %q, want %q", content, "hello world")
	}
}

// BUG #7 (Mistake #91): Multiple temp files with defer in a loop.
func TestMultipleFiles(t *testing.T) {
	files := []struct {
		content  string
		expected string
	}{
		{"first file", "first file"},
		{"second file", "second file"},
	}

	for _, f := range files {
		path, err := WriteTempFile(f.content)
		if err != nil {
			t.Fatalf("WriteTempFile failed: %v", err)
		}
		// BUG: Using defer in a loop - all cleanups happen at function exit
		defer os.Remove(path)

		content, err := ReadFileContent(path)
		if err != nil {
			t.Fatalf("ReadFileContent failed: %v", err)
		}
		if content != f.expected {
			t.Errorf("content = %q, want %q", content, f.expected)
		}
	}
}
