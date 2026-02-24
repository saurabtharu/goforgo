// sync_primitives.go - SOLUTION
// Learn correct usage of sync primitives and avoid subtle concurrency bugs
//
// Mistakes #71-74 from "100 Go Mistakes":
// #71: WaitGroup pitfalls - calling Add inside goroutine or wrong counter
// #72: sync.Cond misuse - not using the associated mutex correctly
// #73: Not collecting errors from goroutines (errgroup-like pattern)
// #74: Copying sync types instead of passing pointers

package main

import (
	"fmt"
	"sort"
	"sync"
)

// === Part 1: WaitGroup Pitfalls (Mistake #71) ===

// FIX: wg.Add(1) is called BEFORE the goroutine is launched
func processJobs(jobs []string) []string {
	var wg sync.WaitGroup
	ch := make(chan string, len(jobs))

	for _, job := range jobs {
		wg.Add(1) // FIX: Add before go statement
		go func(j string) {
			defer wg.Done()
			ch <- fmt.Sprintf("done: %s", j)
		}(job)
	}

	wg.Wait()
	close(ch)

	var results []string
	for r := range ch {
		results = append(results, r)
	}
	sort.Strings(results)
	return results
}

// === Part 2: sync.Cond Misuse (Mistake #72) ===

type WorkQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	items []string
	done  bool
}

func NewWorkQueue() *WorkQueue {
	wq := &WorkQueue{}
	wq.cond = sync.NewCond(&wq.mu)
	return wq
}

// Enqueue adds an item and notifies ONE waiting consumer.
func (wq *WorkQueue) Enqueue(item string) {
	wq.mu.Lock()
	wq.items = append(wq.items, item)
	wq.mu.Unlock()
	wq.cond.Signal()
}

// FIX: Use Broadcast to wake ALL waiting consumers when closing
func (wq *WorkQueue) Close() {
	wq.mu.Lock()
	wq.done = true
	wq.mu.Unlock()
	wq.cond.Broadcast() // FIX: Wake all consumers, not just one
}

// FIX: Use for loop instead of if to re-check condition after wakeup
func (wq *WorkQueue) Dequeue() (string, bool) {
	wq.mu.Lock()
	defer wq.mu.Unlock()

	// FIX: for loop handles spurious wakeups and stolen items
	for len(wq.items) == 0 && !wq.done {
		wq.cond.Wait()
	}

	if len(wq.items) > 0 {
		item := wq.items[0]
		wq.items = wq.items[1:]
		return item, true
	}
	return "", false
}

// === Part 3: Losing Goroutine Errors (Mistake #73) ===

type task struct {
	name string
	fail bool
}

// FIX: Use a buffered error channel to capture the first error
func runTasks(tasks []task) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks)) // FIX: Channel to collect errors

	for _, t := range tasks {
		wg.Add(1)
		go func(t task) {
			defer wg.Done()
			if t.fail {
				errCh <- fmt.Errorf("task %s failed", t.name)
			} else {
				fmt.Printf("task %s: succeeded\n", t.name)
			}
		}(t)
	}

	wg.Wait()
	close(errCh)

	// FIX: Return the first error if any
	for err := range errCh {
		return err
	}
	return nil
}

// === Part 4: Copying Sync Types (Mistake #74) ===

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewSafeMap() *SafeMap { // FIX: Return pointer
	return &SafeMap{data: make(map[string]int)}
}

func (m *SafeMap) Set(key string, val int) {
	m.mu.Lock()
	m.data[key] = val
	m.mu.Unlock()
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

// FIX: Accept *SafeMap (pointer) so the mutex is not copied
func populateMap(m *SafeMap, n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Set(fmt.Sprintf("key-%d", val), val)
		}(i)
	}
	wg.Wait()
}

func main() {
	fmt.Println("=== WaitGroup Pitfalls (Mistake #71) ===")

	jobs := []string{"build", "test", "deploy"}
	results := processJobs(jobs)
	fmt.Printf("completed: %v\n", results)

	expected := []string{"done: build", "done: deploy", "done: test"}
	if len(results) == len(expected) {
		allMatch := true
		for i := range expected {
			if results[i] != expected[i] {
				allMatch = false
				break
			}
		}
		if allMatch {
			fmt.Println("PASS: All jobs completed")
		} else {
			fmt.Println("FAIL: Results don't match expected")
		}
	} else {
		fmt.Printf("FAIL: Expected %d results, got %d\n", len(expected), len(results))
	}

	fmt.Println()
	fmt.Println("=== sync.Cond Misuse (Mistake #72) ===")

	wq := NewWorkQueue()
	var consumed []string
	var consumeMu sync.Mutex
	var consumerWg sync.WaitGroup

	// Start 3 consumers
	for i := 0; i < 3; i++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			for {
				item, ok := wq.Dequeue()
				if !ok {
					return
				}
				consumeMu.Lock()
				consumed = append(consumed, fmt.Sprintf("c%d:%s", id, item))
				consumeMu.Unlock()
			}
		}(i)
	}

	// Produce items
	for _, item := range []string{"a", "b", "c"} {
		wq.Enqueue(item)
	}
	wq.Close()

	consumerWg.Wait()
	sort.Strings(consumed)
	fmt.Printf("consumed: %v\n", consumed)
	if len(consumed) == 3 {
		fmt.Println("PASS: All items consumed")
	} else {
		fmt.Printf("FAIL: Expected 3 items consumed, got %d\n", len(consumed))
	}

	fmt.Println()
	fmt.Println("=== Losing Goroutine Errors (Mistake #73) ===")

	tasks := []task{
		{name: "fetch-data", fail: false},
		{name: "validate", fail: true},
		{name: "transform", fail: false},
	}

	err := runTasks(tasks)
	if err != nil {
		fmt.Printf("correctly caught error: %v\n", err)
		fmt.Println("PASS: Error propagated from goroutine")
	} else {
		fmt.Println("FAIL: Error was lost - runTasks returned nil")
	}

	fmt.Println()
	fmt.Println("=== Copying Sync Types (Mistake #74) ===")

	sm := NewSafeMap()
	populateMap(sm, 10)

	count := 0
	for i := 0; i < 10; i++ {
		if _, ok := sm.Get(fmt.Sprintf("key-%d", i)); ok {
			count++
		}
	}

	fmt.Printf("found %d/10 keys in original map\n", count)
	if count == 10 {
		fmt.Println("PASS: All writes went to the original map")
	} else {
		fmt.Println("FAIL: SafeMap was copied - writes went to the copy")
	}
}
