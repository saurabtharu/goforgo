// functional_options.go
// Refactor a constructor with too many parameters using the functional options pattern
//
// The functional options pattern replaces long parameter lists with
// self-documenting option functions. Instead of:
//   NewServer("api", ":8080", 30, true, 100, "v2")
// You write:
//   NewServer("api", WithAddress(":8080"), WithTimeout(30), WithTLS(true))
//
// Refactor the HTTPClient constructor to use functional options.

package main

import (
	"fmt"
	"time"
)

// HTTPClient makes HTTP requests with various configuration options.
type HTTPClient struct {
	baseURL        string
	timeout        time.Duration
	maxRetries     int
	retryDelay     time.Duration
	followRedirect bool
	maxRedirects   int
	userAgent      string
	headers        map[string]string
}

// BUG: This constructor has too many parameters. Callers must remember
// the order and provide values for everything, even defaults.
//
// TODO: Refactor to use the functional options pattern:
// 1. Define an Option type: type Option func(*HTTPClient)
// 2. Create With... functions for each option
// 3. Change NewHTTPClient to accept (baseURL string, opts ...Option)
// 4. Set sensible defaults, then apply options
func NewHTTPClient(
	baseURL string,
	timeout time.Duration,
	maxRetries int,
	retryDelay time.Duration,
	followRedirect bool,
	maxRedirects int,
	userAgent string,
	headers map[string]string,
) *HTTPClient {
	return &HTTPClient{
		baseURL:        baseURL,
		timeout:        timeout,
		maxRetries:     maxRetries,
		retryDelay:     retryDelay,
		followRedirect: followRedirect,
		maxRedirects:   maxRedirects,
		userAgent:      userAgent,
		headers:        headers,
	}
}

func (c *HTTPClient) Describe() {
	fmt.Printf("HTTPClient Configuration:\n")
	fmt.Printf("  Base URL:         %s\n", c.baseURL)
	fmt.Printf("  Timeout:          %s\n", c.timeout)
	fmt.Printf("  Max Retries:      %d\n", c.maxRetries)
	fmt.Printf("  Retry Delay:      %s\n", c.retryDelay)
	fmt.Printf("  Follow Redirect:  %v\n", c.followRedirect)
	fmt.Printf("  Max Redirects:    %d\n", c.maxRedirects)
	fmt.Printf("  User Agent:       %s\n", c.userAgent)
	if len(c.headers) > 0 {
		fmt.Println("  Custom Headers:")
		for k, v := range c.headers {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
}

func main() {
	fmt.Println("=== Functional Options Pattern ===")

	// BUG: This call is hard to read. What does 30, 3, true, 10 mean?
	// TODO: After refactoring, create the clients using functional options.

	// Client 1: Simple client with mostly defaults
	// After refactoring, this should be:
	//   client1 := NewHTTPClient("https://api.example.com")
	fmt.Println("\n--- Client 1: Defaults ---")
	client1 := NewHTTPClient(
		"https://api.example.com",
		30*time.Second,  // timeout
		3,               // maxRetries
		1*time.Second,   // retryDelay
		true,            // followRedirect
		10,              // maxRedirects
		"GoForGo/1.0",  // userAgent
		nil,             // headers
	)
	client1.Describe()

	// Client 2: Custom timeout and retries
	// After refactoring, this should be:
	//   client2 := NewHTTPClient("https://slow.example.com",
	//       WithTimeout(60*time.Second),
	//       WithMaxRetries(5),
	//       WithRetryDelay(2*time.Second),
	//   )
	fmt.Println("\n--- Client 2: Custom timeout ---")
	client2 := NewHTTPClient(
		"https://slow.example.com",
		60*time.Second,
		5,
		2*time.Second,
		true,
		10,
		"GoForGo/1.0",
		nil,
	)
	client2.Describe()

	// Client 3: Full customization with headers
	// After refactoring, this should be:
	//   client3 := NewHTTPClient("https://internal.example.com",
	//       WithTimeout(10*time.Second),
	//       WithNoRedirects(),
	//       WithUserAgent("InternalBot/2.0"),
	//       WithHeader("Authorization", "Bearer token123"),
	//       WithHeader("X-Request-ID", "abc-456"),
	//   )
	fmt.Println("\n--- Client 3: Full custom ---")
	client3 := NewHTTPClient(
		"https://internal.example.com",
		10*time.Second,
		3,
		1*time.Second,
		false,
		0,
		"InternalBot/2.0",
		map[string]string{
			"Authorization": "Bearer token123",
			"X-Request-ID":  "abc-456",
		},
	)
	client3.Describe()

	fmt.Println("\nFunctional options refactoring complete!")
}
