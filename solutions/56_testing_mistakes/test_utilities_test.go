// test_utilities_test.go
// Solution for Mistakes #88, #90, #91
//
// Fixed:
// 1. HTTP tests use httptest.NewRecorder and http.NewRequest (Mistake #88)
// 2. Helper functions call t.Helper() for accurate error reporting (Mistake #90)
// 3. Resource cleanup uses t.Cleanup instead of defer (Mistake #91)

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// FIXED: t.Helper() added so error messages point to the calling test.
func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}
}

// FIXED: t.Helper() added.
func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if len(body) == 0 || len(substr) == 0 {
		if substr != "" {
			t.Errorf("body is empty, expected to contain %q", substr)
		}
		return
	}
	for i := 0; i <= len(body)-len(substr); i++ {
		if body[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("body %q does not contain %q", body, substr)
}

// FIXED: Uses httptest.NewRecorder and http.NewRequest to properly
// test the HTTP handler including status codes, headers, and body.
func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	HealthHandler(rec, req)

	assertStatusCode(t, rec.Code, http.StatusOK)

	// Verify Content-Type header
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	// Verify response body
	body, _ := io.ReadAll(rec.Body)
	assertContains(t, string(body), `"status":"ok"`)
}

// FIXED: Uses httptest for both success and error cases.
func TestUserHandler(t *testing.T) {
	t.Run("valid id", func(t *testing.T) {
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
		if user.Name != "Alice" {
			t.Errorf("user.Name = %q, want %q", user.Name, "Alice")
		}
	})

	t.Run("missing id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user", nil)

		UserHandler(rec, req)

		assertStatusCode(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("invalid id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user?id=abc", nil)

		UserHandler(rec, req)

		assertStatusCode(t, rec.Code, http.StatusBadRequest)
	})
}

// FIXED: Uses httptest for handler testing.
func TestGreetHandler(t *testing.T) {
	t.Run("default name", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/greet", nil)

		GreetHandler(rec, req)

		assertStatusCode(t, rec.Code, http.StatusOK)
		body, _ := io.ReadAll(rec.Body)
		if string(body) != "Hello, World!" {
			t.Errorf("body = %q, want %q", string(body), "Hello, World!")
		}
	})

	t.Run("custom name", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/greet?name=Gopher", nil)

		GreetHandler(rec, req)

		assertStatusCode(t, rec.Code, http.StatusOK)
		body, _ := io.ReadAll(rec.Body)
		if string(body) != "Hello, Gopher!" {
			t.Errorf("body = %q, want %q", string(body), "Hello, Gopher!")
		}
	})
}

// FIXED: Uses t.Cleanup instead of defer for resource cleanup.
func TestWriteAndReadFile(t *testing.T) {
	path, err := WriteTempFile("hello world")
	if err != nil {
		t.Fatalf("WriteTempFile failed: %v", err)
	}
	// FIXED: t.Cleanup is tied to the test lifecycle and runs even on panic
	t.Cleanup(func() { os.Remove(path) })

	content, err := ReadFileContent(path)
	if err != nil {
		t.Fatalf("ReadFileContent failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("content = %q, want %q", content, "hello world")
	}
}

// FIXED: Uses t.Cleanup for each temp file in the loop.
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
		// FIXED: t.Cleanup per iteration ensures each file gets cleaned up
		t.Cleanup(func() { os.Remove(path) })

		content, err := ReadFileContent(path)
		if err != nil {
			t.Fatalf("ReadFileContent failed: %v", err)
		}
		if content != f.expected {
			t.Errorf("content = %q, want %q", content, f.expected)
		}
	}
}
