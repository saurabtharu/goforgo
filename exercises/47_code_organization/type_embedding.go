// type_embedding.go
// Fix unintended method promotion through struct embedding
//
// When you embed a type in a struct, ALL of its exported methods get
// promoted to the outer struct. This can accidentally expose internal
// implementation details that consumers shouldn't use directly.
//
// This program has a Logger embedded in a Server, which means the
// server accidentally exposes all logging internals. Fix it by using
// a named field instead of embedding, and add wrapper methods where needed.

package main

import (
	"fmt"
	"strings"
	"time"
)

// Logger handles structured log output.
type Logger struct {
	prefix  string
	entries []string
	level   int // 0=debug, 1=info, 2=warn, 3=error
}

func NewLogger(prefix string, level int) *Logger {
	return &Logger{prefix: prefix, level: level}
}

func (l *Logger) Debug(msg string) {
	if l.level <= 0 {
		l.log("DEBUG", msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.level <= 1 {
		l.log("INFO", msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level <= 2 {
		l.log("WARN", msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.level <= 3 {
		l.log("ERROR", msg)
	}
}

// SetLevel changes the log level. This is an internal detail.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// GetEntries returns all log entries. This is for testing/debugging only.
func (l *Logger) GetEntries() []string {
	return l.entries
}

// Reset clears all log entries. Dangerous if called accidentally.
func (l *Logger) Reset() {
	l.entries = nil
}

func (l *Logger) log(level, msg string) {
	entry := fmt.Sprintf("[%s] %s: %s - %s",
		time.Now().Format("15:04:05"), l.prefix, level, msg)
	l.entries = append(l.entries, entry)
	fmt.Println(entry)
}

// Server represents an HTTP server with logging capabilities.
// BUG: Embedding Logger promotes ALL logger methods to Server.
// This means callers can do server.SetLevel(), server.Reset(), etc.
// which exposes internal logger details.
//
// TODO: Change the embedding to a named field (e.g., logger *Logger)
// and only expose the logging methods that Server consumers need:
// - Log(msg string) - calls logger.Info()
// That's it. SetLevel, Reset, GetEntries should NOT be accessible
// through Server.
type Server struct {
	*Logger
	name    string
	address string
	routes  map[string]string
}

func NewServer(name, address string) *Server {
	return &Server{
		Logger:  NewLogger(name, 1), // info level
		name:    name,
		address: address,
		routes:  make(map[string]string),
	}
}

func (s *Server) AddRoute(path, handler string) {
	s.routes[path] = handler
	// BUG: This works because Logger is embedded, promoting Info().
	// After fixing, this should call s.Log() or s.logger.Info() instead.
	s.Info(fmt.Sprintf("Route added: %s -> %s", path, handler))
}

func (s *Server) Start() {
	s.Info(fmt.Sprintf("Starting server %s on %s", s.name, s.address))
	for path, handler := range s.routes {
		s.Debug(fmt.Sprintf("  Registered: %s -> %s", path, handler))
	}
	s.Info("Server started successfully")
}

func (s *Server) HandleRequest(path string) {
	handler, ok := s.routes[path]
	if !ok {
		s.Warn(fmt.Sprintf("No handler for path: %s", path))
		return
	}
	s.Info(fmt.Sprintf("Handling %s with %s", path, handler))
}

func main() {
	fmt.Println("=== Type Embedding ===")

	server := NewServer("api-server", ":8080")

	// Add some routes
	server.AddRoute("/health", "healthCheck")
	server.AddRoute("/users", "listUsers")
	server.AddRoute("/orders", "listOrders")

	// Start the server
	server.Start()

	// Handle some requests
	fmt.Println()
	server.HandleRequest("/users")
	server.HandleRequest("/missing")

	// BUG: These calls should NOT be possible on a Server!
	// After fixing, these lines will cause compile errors.
	// Comment them out or remove them after refactoring.
	fmt.Println("\n--- Exposed internal methods (BUG) ---")
	server.SetLevel(0) // Should not be accessible!
	server.Reset()     // Should not be accessible!
	entries := server.GetEntries()
	fmt.Println("Entries after reset:", len(entries))

	// After fix, the server should only expose a Log method:
	fmt.Println("\n--- After fixing ---")
	server.Log("This uses the clean public API")

	fmt.Println("\nEmbedding bugs fixed!")
}

// Log is the only logging method that Server should expose.
// TODO: This is already defined, but it uses the embedded logger directly.
// After switching to a named field, update this to use the field.
func (s *Server) Log(msg string) {
	s.Info(msg)
}

// printSummary shows how many log entries exist.
func printSummary(entries []string) {
	fmt.Printf("Total log entries: %d\n", len(entries))
	for _, e := range entries {
		short := e
		if len(short) > 80 {
			short = short[:80] + "..."
		}
		// Just show the message part
		parts := strings.SplitN(short, " - ", 2)
		if len(parts) == 2 {
			fmt.Println(" ", parts[1])
		}
	}
}
