// functional_options.go - SOLUTION
// Refactored to use the functional options pattern.

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

// Option is a function that configures an HTTPClient.
type Option func(*HTTPClient)

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *HTTPClient) {
		c.timeout = d
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) Option {
	return func(c *HTTPClient) {
		c.maxRetries = n
	}
}

// WithRetryDelay sets the delay between retries.
func WithRetryDelay(d time.Duration) Option {
	return func(c *HTTPClient) {
		c.retryDelay = d
	}
}

// WithNoRedirects disables following redirects.
func WithNoRedirects() Option {
	return func(c *HTTPClient) {
		c.followRedirect = false
		c.maxRedirects = 0
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *HTTPClient) {
		c.userAgent = ua
	}
}

// WithHeader adds a custom header.
func WithHeader(key, value string) Option {
	return func(c *HTTPClient) {
		if c.headers == nil {
			c.headers = make(map[string]string)
		}
		c.headers[key] = value
	}
}

// Fixed: Accepts variadic options with sensible defaults.
func NewHTTPClient(baseURL string, opts ...Option) *HTTPClient {
	c := &HTTPClient{
		baseURL:        baseURL,
		timeout:        30 * time.Second,
		maxRetries:     3,
		retryDelay:     1 * time.Second,
		followRedirect: true,
		maxRedirects:   10,
		userAgent:      "GoForGo/1.0",
		headers:        make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
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
	} else {
		fmt.Println("  Custom Headers:")
	}
}

func main() {
	fmt.Println("=== Functional Options Pattern ===")

	// Client 1: Simple client with all defaults
	fmt.Println("\n--- Client 1: Defaults ---")
	client1 := NewHTTPClient("https://api.example.com")
	client1.Describe()

	// Client 2: Custom timeout and retries - only specify what differs
	fmt.Println("\n--- Client 2: Custom timeout ---")
	client2 := NewHTTPClient("https://slow.example.com",
		WithTimeout(60*time.Second),
		WithMaxRetries(5),
		WithRetryDelay(2*time.Second),
	)
	client2.Describe()

	// Client 3: Full customization - self-documenting option names
	fmt.Println("\n--- Client 3: Full custom ---")
	client3 := NewHTTPClient("https://internal.example.com",
		WithTimeout(10*time.Second),
		WithNoRedirects(),
		WithUserAgent("InternalBot/2.0"),
		WithHeader("Authorization", "Bearer token123"),
		WithHeader("X-Request-ID", "abc-456"),
	)
	client3.Describe()

	fmt.Println("\nFunctional options refactoring complete!")
}
