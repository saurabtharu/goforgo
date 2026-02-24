// data_race_patterns.go - SOLUTION
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

// FIX: updateUser holds the mutex while modifying fields
func updateUser(u *User, name string, age int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Name = name
	u.Age = age
}

// FIX: readUser holds the mutex while reading fields
func readUser(u *User) string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.String()
}

// === Part 2: Append Race (Mistake #69) ===

// FIX: Use a channel to collect results instead of concurrent append
func collectResults(n int) []string {
	ch := make(chan string, n)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// FIX: Send through channel instead of appending to shared slice
			ch <- fmt.Sprintf("result-%d", id)
		}(i)
	}

	wg.Wait()
	close(ch)

	results := make([]string, 0, n)
	for r := range ch {
		results = append(results, r)
	}
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

// FIX: Accept *Counter (pointer) so the mutex is not copied
func incrementMany(c *Counter, n int) {
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

	counter := &Counter{} // FIX: Use pointer
	incrementMany(counter, 100)

	fmt.Printf("counter value: %d (expected 100)\n", counter.Value())
	if counter.Value() == 100 {
		fmt.Println("PASS: Counter incremented correctly")
	} else {
		fmt.Println("FAIL: Counter passed by value - mutex was copied!")
	}
}
