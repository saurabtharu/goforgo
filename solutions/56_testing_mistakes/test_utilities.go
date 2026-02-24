// test_utilities.go
// Solution for Mistakes #88, #90, #91
//
// The source code is the same - the fixes are in the test file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// User represents a simple user record.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// HealthHandler returns a health check response.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// UserHandler returns a user by ID from the query string.
func UserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	user := User{ID: id, Name: "Alice", Age: 30}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GreetHandler returns a greeting message.
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, %s!", name)
}

// WriteTempFile creates a temporary file with the given content.
// The caller is responsible for removing the file.
func WriteTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

// ReadFileContent reads and returns the content of a file.
func ReadFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	fmt.Println("Run 'go test -v' to execute the tests.")
}
