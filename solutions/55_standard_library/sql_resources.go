package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// 100 Go Mistakes #78, #79: SQL resource management
//
// Solution: Always close rows with defer, check rows.Err(), and configure pool settings.

func main() {
	fmt.Println("=== SQL Resource Mistakes ===")
	fmt.Println()

	testRowsClosing()
	fmt.Println()

	testRowsErr()
	fmt.Println()

	testConnectionPool()
	fmt.Println()

	fmt.Println("All SQL resource checks passed!")
}

// --- Mock Database Types ---

type MockDB struct {
	mu              sync.Mutex
	openConnections int
	maxOpen         int
	maxIdle         int
	connMaxLifetime time.Duration
	queryCount      int
}

type MockRows struct {
	db      *MockDB
	data    []string
	index   int
	closed  bool
	scanErr error
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (db *MockDB) SetMaxOpenConns(n int)                { db.maxOpen = n }
func (db *MockDB) SetMaxIdleConns(n int)                { db.maxIdle = n }
func (db *MockDB) SetConnMaxLifetime(d time.Duration)   { db.connMaxLifetime = d }

func (db *MockDB) Query(query string) *MockRows {
	db.mu.Lock()
	db.openConnections++
	db.queryCount++
	db.mu.Unlock()
	return &MockRows{
		db:   db,
		data: []string{"Alice", "Bob", "Charlie"},
	}
}

func (db *MockDB) QueryWithError(query string) *MockRows {
	db.mu.Lock()
	db.openConnections++
	db.queryCount++
	db.mu.Unlock()
	return &MockRows{
		db:      db,
		data:    []string{"Alice", "Bob", "Charlie"},
		scanErr: fmt.Errorf("connection reset during scan"),
	}
}

func (r *MockRows) Next() bool {
	if r.closed {
		return false
	}
	if r.scanErr != nil && r.index == 2 {
		return false
	}
	return r.index < len(r.data)
}

func (r *MockRows) Scan(dest *string) error {
	*dest = r.data[r.index]
	r.index++
	return nil
}

func (r *MockRows) Err() error {
	return r.scanErr
}

func (r *MockRows) Close() error {
	if !r.closed {
		r.closed = true
		r.db.mu.Lock()
		r.db.openConnections--
		r.db.mu.Unlock()
	}
	return nil
}

// --- Test 1: Forgetting to close rows ---

func testRowsClosing() {
	fmt.Println("--- Rows Closing ---")

	db := NewMockDB()
	names := queryUsers(db)
	fmt.Printf("Got users: %v\n", names)

	if db.openConnections == 0 {
		fmt.Println("PASS: all connections returned to pool")
	} else {
		fmt.Printf("FAIL: %d connections still open (leaked!)\n", db.openConnections)
	}
}

// FIXED: Added defer rows.Close() immediately after Query
func queryUsers(db *MockDB) []string {
	rows := db.Query("SELECT name FROM users")
	defer rows.Close() // FIXED: Always close rows!

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}
	return names
}

// --- Test 2: Not checking rows.Err() ---

func testRowsErr() {
	fmt.Println("--- Rows Error Check ---")

	db := NewMockDB()
	names, err := queryUsersWithErrorCheck(db)
	fmt.Printf("Got users: %v\n", names)

	if err != nil {
		fmt.Printf("PASS: error detected: %v\n", err)
	} else {
		fmt.Println("FAIL: should have caught the iteration error")
	}
}

// FIXED: Check rows.Err() after the iteration loop
func queryUsersWithErrorCheck(db *MockDB) ([]string, error) {
	rows := db.QueryWithError("SELECT name FROM users")
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	// FIXED: Check for errors that stopped the iteration
	if err := rows.Err(); err != nil {
		return names, err
	}

	return names, nil
}

// --- Test 3: Connection pool configuration ---

func testConnectionPool() {
	fmt.Println("--- Connection Pool ---")

	db := NewMockDB()

	// FIXED: Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	var results []string
	results = append(results, checkPool("MaxOpenConns", db.maxOpen, 25))
	results = append(results, checkPool("MaxIdleConns", db.maxIdle, 10))

	if db.connMaxLifetime == 5*time.Minute {
		fmt.Println("PASS: ConnMaxLifetime set to 5m")
		results = append(results, "pass")
	} else {
		fmt.Printf("FAIL: ConnMaxLifetime is %v, expected 5m\n", db.connMaxLifetime)
		results = append(results, "fail")
	}

	allPass := true
	for _, r := range results {
		if r == "fail" {
			allPass = false
		}
	}
	if allPass {
		fmt.Println("PASS: connection pool properly configured")
	}
}

func checkPool(name string, got, want int) string {
	if got == want {
		fmt.Printf("PASS: %s = %d\n", name, got)
		return "pass"
	}
	status := strings.Builder{}
	status.WriteString(fmt.Sprintf("FAIL: %s = %d, expected %d", name, got, want))
	fmt.Println(status.String())
	return "fail"
}
