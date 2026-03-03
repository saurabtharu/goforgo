// data_race_patterns.go
// Learn to identify and fix common data race patterns in Go
//
// Mistakes #68-70 from "100 Go Mistakes":
// #68: String formatting race - concurrent reads via fmt while another goroutine writes
// #69: Append race - multiple goroutines appending to the same slice
// #70: Mutex misuse - copying a mutex or holding it with wrong granularity

package main

import (
	"fmt"
	"sort"
	"sync"
)

// === Part 1: String Formatting Race (Mistake #68) ===

type User struct {
	mu   sync.Mutex
	Name string
	Age  int
}

func (u *User) String() string {
	return fmt.Sprintf("%s (age %d)", u.Name, u.Age)
}

// updateUser modifies the user's fields.
func updateUser(u *User, name string, age int) {
	u.Name = name
	u.Age = age
}

// readUser formats the user as a string.
//
// BUG: This calls u.String() which reads Name and Age without holding the lock.
// If another goroutine is calling updateUser concurrently, this is a data race.
//
// TODO: Lock the mutex before reading user fields. Both readUser and updateUser
// need to coordinate via the mutex.
func readUser(u *User) string {
	return u.String()
}

// === Part 2: Append Race (Mistake #69) ===

// collectResults has multiple goroutines appending to a shared slice.
//
// BUG: Concurrent append to the same slice is a data race. Slice headers
// (pointer, length, capacity) can be corrupted by concurrent writes.
//
// TODO: Instead of appending to a shared slice, have each goroutine send
// its result through a channel. Collect results from the channel after
// all goroutines finish.
func collectResults(n int) []string {
	results := make([]string, 0, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// BUG: Concurrent append to shared slice - data race!
			results = append(results, fmt.Sprintf("result-%d", id))
		}(i)
	}

	wg.Wait()
	sort.Strings(results)
	return results
}

// === Part 3: Mutex Misuse (Mistake #70) ===

type Counter struct {
	mu    sync.Mutex
	value int
}

// Increment safely increments the counter.
func (c *Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// Value safely reads the counter value.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// incrementMany runs n goroutines that each increment the counter.
//
// BUG: The counter is passed by VALUE, which copies the mutex.
// A copied mutex doesn't protect the original - each copy has its own
// independent lock state.
//
// TODO: Pass the counter by pointer (*Counter) instead of by value.
func incrementMany(c Counter, n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()
}

func main() {
	fmt.Println("=== String Formatting Race (Mistake #68) ===")

	user := &User{Name: "Alice", Age: 30}
	var wg sync.WaitGroup
	results := make([]string, 2)

	// Writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		updateUser(user, "Bob", 25)
		results[0] = "update done"
	}()

	// Reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		snapshot := readUser(user)
		results[1] = fmt.Sprintf("read: %s", snapshot)
	}()

	wg.Wait()
	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Println("(with -race flag, the buggy version would report a data race)")

	fmt.Println()
	fmt.Println("=== Append Race (Mistake #69) ===")

	collected := collectResults(5)
	fmt.Printf("collected %d results: %v\n", len(collected), collected)

	expected := []string{"result-0", "result-1", "result-2", "result-3", "result-4"}
	if len(collected) == len(expected) {
		allMatch := true
		for i := range expected {
			if collected[i] != expected[i] {
				allMatch = false
				break
			}
		}
		if allMatch {
			fmt.Println("PASS: All results collected safely")
		} else {
			fmt.Println("FAIL: Results corrupted by data race")
		}
	} else {
		fmt.Printf("FAIL: Expected %d results, got %d (data race lost some)\n", len(expected), len(collected))
	}

	fmt.Println()
	fmt.Println("=== Mutex Misuse (Mistake #70) ===")

	counter := Counter{}
	incrementMany(counter, 100)

	fmt.Printf("counter value: %d (expected 100)\n", counter.Value())
	if counter.Value() == 100 {
		fmt.Println("PASS: Counter incremented correctly")
	} else {
		fmt.Println("FAIL: Counter passed by value - mutex was copied!")
	}
}
