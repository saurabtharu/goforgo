// type_embedding.go - SOLUTION
// Fixed: Named field instead of embedding. Only Log() is exposed on Server.

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
	level   int
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

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) GetEntries() []string {
	return l.entries
}

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
// Fixed: Named field 'logger' instead of embedding *Logger.
// Only the Log() method is exposed to Server consumers.
type Server struct {
	logger  *Logger // Named field - no method promotion
	name    string
	address string
	routes  map[string]string
}

func NewServer(name, address string) *Server {
	return &Server{
		logger:  NewLogger(name, 1),
		name:    name,
		address: address,
		routes:  make(map[string]string),
	}
}

// Log is the only logging method exposed by Server.
func (s *Server) Log(msg string) {
	s.logger.Info(msg)
}

// Fixed: Uses s.logger.Info() through the named field.
func (s *Server) AddRoute(path, handler string) {
	s.routes[path] = handler
	s.logger.Info(fmt.Sprintf("Route added: %s -> %s", path, handler))
}

func (s *Server) Start() {
	s.logger.Info(fmt.Sprintf("Starting server %s on %s", s.name, s.address))
	for path, handler := range s.routes {
		s.logger.Debug(fmt.Sprintf("  Registered: %s -> %s", path, handler))
	}
	s.logger.Info("Server started successfully")
}

func (s *Server) HandleRequest(path string) {
	handler, ok := s.routes[path]
	if !ok {
		s.logger.Warn(fmt.Sprintf("No handler for path: %s", path))
		return
	}
	s.logger.Info(fmt.Sprintf("Handling %s with %s", path, handler))
}

func main() {
	fmt.Println("=== Type Embedding ===")

	server := NewServer("api-server", ":8080")

	server.AddRoute("/health", "healthCheck")
	server.AddRoute("/users", "listUsers")
	server.AddRoute("/orders", "listOrders")

	server.Start()

	fmt.Println()
	server.HandleRequest("/users")
	server.HandleRequest("/missing")

	// Fixed: These lines are removed. SetLevel, Reset, GetEntries
	// are no longer accessible on Server - only on Logger directly.

	fmt.Println("\n--- After fixing ---")
	server.Log("This uses the clean public API")

	fmt.Println("\nEmbedding bugs fixed!")
}

func printSummary(entries []string) {
	fmt.Printf("Total log entries: %d\n", len(entries))
	for _, e := range entries {
		short := e
		if len(short) > 80 {
			short = short[:80] + "..."
		}
		parts := strings.SplitN(short, " - ", 2)
		if len(parts) == 2 {
			fmt.Println(" ", parts[1])
		}
	}
}
