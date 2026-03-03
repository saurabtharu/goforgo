package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// 100 Go Mistakes #78, #79: SQL resource management
//
// This exercise simulates database patterns without requiring a real database.
// It covers two critical mistakes:
// 1. Not closing rows / not checking rows.Err() after iteration
// 2. Not configuring connection pool settings (MaxOpen, MaxIdle, ConnMaxLifetime)
//
// We use mock types that track resource usage to demonstrate these patterns.

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
	scanErr error // Simulates an error partway through iteration
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
		return false // Stop early due to error
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

	// BUG: This function queries but never closes the rows!
	// When you forget to call rows.Close(), the underlying database connection
	// is never returned to the pool, causing a connection leak.
	//
	// TODO: Use defer rows.Close() immediately after the Query call.
	// Always close rows even if you return early from the function.

	names := queryUsers(db)
	fmt.Printf("Got users: %v\n", names)

	if db.openConnections == 0 {
		fmt.Println("PASS: all connections returned to pool")
	} else {
		fmt.Printf("FAIL: %d connections still open (leaked!)\n", db.openConnections)
	}
}

func queryUsers(db *MockDB) []string {
	rows := db.Query("SELECT name FROM users")
	// BUG: Missing rows.Close()!
	// TODO: Add defer rows.Close() here

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

	// BUG: After iterating with rows.Next(), you must check rows.Err()
	// because Next() can stop due to an error, not just end of data.
	// Ignoring this means you silently process incomplete results.
	//
	// TODO: Check rows.Err() after the loop and report any error found.

	names, err := queryUsersWithErrorCheck(db)
	fmt.Printf("Got users: %v\n", names)

	if err != nil {
		fmt.Printf("PASS: error detected: %v\n", err)
	} else {
		fmt.Println("FAIL: should have caught the iteration error")
	}
}

func queryUsersWithErrorCheck(db *MockDB) ([]string, error) {
	rows := db.QueryWithError("SELECT name FROM users")
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	// BUG: Not checking rows.Err()!
	// TODO: Check rows.Err() here and return it if non-nil.

	return names, nil
}

// --- Test 3: Connection pool configuration ---

func testConnectionPool() {
	fmt.Println("--- Connection Pool ---")

	// BUG: Using the database without configuring connection pool settings.
	// The defaults are:
	// - MaxOpenConns: 0 (unlimited!) - can overwhelm the database
	// - MaxIdleConns: 2 (too low for most apps) - constant reconnection overhead
	// - ConnMaxLifetime: 0 (connections live forever) - stale connections accumulate
	//
	// TODO: Set reasonable pool settings:
	// - MaxOpenConns: 25
	// - MaxIdleConns: 10
	// - ConnMaxLifetime: 5 minutes

	db := NewMockDB()

	// Missing pool configuration!
	// TODO: Add db.SetMaxOpenConns(), db.SetMaxIdleConns(), db.SetConnMaxLifetime()

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
