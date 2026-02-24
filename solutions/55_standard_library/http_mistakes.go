package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// 100 Go Mistakes #80, #81: HTTP handler and client/server configuration mistakes
//
// Solution: Always return after http.Error(), configure client/server timeouts.

func main() {
	fmt.Println("=== HTTP Mistakes ===")
	fmt.Println()

	testMissingReturn()
	fmt.Println()

	testHTTPClientConfig()
	fmt.Println()

	testHTTPServerConfig()
	fmt.Println()

	fmt.Println("All HTTP checks passed!")
}

// --- Mistake #80: Missing return after http.Error ---

// FIXED: Added return after each http.Error() call
func buggyHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "missing API key", http.StatusUnauthorized)
		return // FIXED: Return after error
	}

	if apiKey != "secret-key-123" {
		http.Error(w, "invalid API key", http.StatusForbidden)
		return // FIXED: Return after error
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"data": "sensitive information", "user": "admin"}`)
}

func testMissingReturn() {
	fmt.Println("--- Missing Return After http.Error ---")

	// Test 1: No API key
	req := httptest.NewRequest("GET", "/api/data", nil)
	rec := httptest.NewRecorder()
	buggyHandler(rec, req)

	body := rec.Body.String()
	fmt.Printf("No-key response: %s\n", strings.TrimSpace(body))

	if strings.Contains(body, "sensitive information") {
		fmt.Println("FAIL: sensitive data leaked in error response!")
	} else {
		fmt.Println("PASS: error response contains only the error")
	}

	// Test 2: Wrong API key
	req2 := httptest.NewRequest("GET", "/api/data", nil)
	req2.Header.Set("X-API-Key", "wrong-key")
	rec2 := httptest.NewRecorder()
	buggyHandler(rec2, req2)

	body2 := rec2.Body.String()
	fmt.Printf("Wrong-key response: %s\n", strings.TrimSpace(body2))

	if strings.Contains(body2, "sensitive information") {
		fmt.Println("FAIL: sensitive data leaked with wrong key!")
	} else {
		fmt.Println("PASS: wrong key only returns forbidden error")
	}

	// Test 3: Valid request
	req3 := httptest.NewRequest("GET", "/api/data", nil)
	req3.Header.Set("X-API-Key", "secret-key-123")
	rec3 := httptest.NewRecorder()
	buggyHandler(rec3, req3)

	body3 := rec3.Body.String()
	if strings.Contains(body3, "sensitive information") && !strings.Contains(body3, "missing") {
		fmt.Println("PASS: valid request returns data correctly")
	} else {
		fmt.Println("FAIL: valid request should return sensitive data")
	}
}

// --- Mistake #81: Default HTTP client and server ---

type ClientConfig struct {
	Timeout             time.Duration
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
}

type ServerConfig struct {
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func testHTTPClientConfig() {
	fmt.Println("--- HTTP Client Config ---")

	client := createHTTPClient()
	config := inspectClient(client)

	fmt.Printf("Client timeout: %v\n", config.Timeout)

	checks := 0
	if config.Timeout == 10*time.Second {
		fmt.Println("PASS: client timeout set to 10s")
		checks++
	} else {
		fmt.Printf("FAIL: client timeout is %v, expected 10s\n", config.Timeout)
	}

	if config.IdleConnTimeout == 90*time.Second {
		fmt.Println("PASS: idle connection timeout set to 90s")
		checks++
	} else {
		fmt.Printf("FAIL: idle conn timeout is %v, expected 90s\n", config.IdleConnTimeout)
	}

	if config.TLSHandshakeTimeout == 5*time.Second {
		fmt.Println("PASS: TLS handshake timeout set to 5s")
		checks++
	} else {
		fmt.Printf("FAIL: TLS handshake timeout is %v, expected 5s\n", config.TLSHandshakeTimeout)
	}

	if checks == 3 {
		fmt.Println("PASS: HTTP client fully configured")
	}
}

// FIXED: Configure the client with proper timeouts and transport settings
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}

func inspectClient(c *http.Client) ClientConfig {
	config := ClientConfig{
		Timeout: c.Timeout,
	}
	if t, ok := c.Transport.(*http.Transport); ok {
		config.IdleConnTimeout = t.IdleConnTimeout
		config.TLSHandshakeTimeout = t.TLSHandshakeTimeout
	}
	return config
}

func testHTTPServerConfig() {
	fmt.Println("--- HTTP Server Config ---")

	server := createHTTPServer()
	config := inspectServer(server)

	checks := 0
	if config.ReadTimeout == 5*time.Second {
		fmt.Println("PASS: read timeout set to 5s")
		checks++
	} else {
		fmt.Printf("FAIL: read timeout is %v, expected 5s\n", config.ReadTimeout)
	}

	if config.WriteTimeout == 10*time.Second {
		fmt.Println("PASS: write timeout set to 10s")
		checks++
	} else {
		fmt.Printf("FAIL: write timeout is %v, expected 10s\n", config.WriteTimeout)
	}

	if config.IdleTimeout == 120*time.Second {
		fmt.Println("PASS: idle timeout set to 120s")
		checks++
	} else {
		fmt.Printf("FAIL: idle timeout is %v, expected 120s\n", config.IdleTimeout)
	}

	if config.ReadHeaderTimeout == 2*time.Second {
		fmt.Println("PASS: read header timeout set to 2s")
		checks++
	} else {
		fmt.Printf("FAIL: read header timeout is %v, expected 2s\n", config.ReadHeaderTimeout)
	}

	if checks == 4 {
		fmt.Println("PASS: HTTP server fully configured")
	}
}

// FIXED: Configure all four timeout fields
func createHTTPServer() *http.Server {
	return &http.Server{
		Addr:              ":8080",
		Handler:           http.DefaultServeMux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
}

func inspectServer(s *http.Server) ServerConfig {
	return ServerConfig{
		ReadTimeout:       s.ReadTimeout,
		WriteTimeout:      s.WriteTimeout,
		IdleTimeout:       s.IdleTimeout,
		ReadHeaderTimeout: s.ReadHeaderTimeout,
	}
}
