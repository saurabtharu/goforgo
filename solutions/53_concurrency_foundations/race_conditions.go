// race_conditions.go - SOLUTION
// Understand data races vs race conditions, and I/O-bound vs CPU-bound workloads.
// Based on 100 Go Mistakes #58 and #59.

package main

import (
	"fmt"
	"sync"
)

// === Part 1: Data Race (FIXED) ===
// Mutex protects all concurrent reads and writes to the counter.

func dataRaceDemo() {
	fmt.Println("--- Part 1: Data Race ---")
	counter := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1_000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("Counter: %d (expected: 5000)\n", counter)
	if counter == 5_000 {
		fmt.Println("OK: No lost updates")
	} else {
		fmt.Println("BUG: Lost updates due to data race!")
	}
}

// === Part 2: Race Condition (FIXED) ===
// The entire check-and-withdraw is now a single critical section.

type Account struct {
	mu      sync.Mutex
	balance int
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// transfer holds the lock for the entire check-and-withdraw sequence,
// making it atomic. No other goroutine can see or modify the balance
// between the check and the withdrawal.
func transfer(from *Account, amount int, id int, results chan<- string) {
	from.mu.Lock()
	if from.balance >= amount {
		from.balance -= amount
		remaining := from.balance
		from.mu.Unlock()
		results <- fmt.Sprintf("Transfer %d: withdrew %d, remaining: %d", id, amount, remaining)
	} else {
		bal := from.balance
		from.mu.Unlock()
		results <- fmt.Sprintf("Transfer %d: insufficient funds (balance: %d)", id, bal)
	}
}

func raceConditionDemo() {
	fmt.Println("\n--- Part 2: Race Condition ---")
	account := &Account{balance: 100}
	results := make(chan string, 3)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			transfer(account, 60, id, results)
		}(i)
	}
	wg.Wait()
	close(results)

	for msg := range results {
		fmt.Println(msg)
	}
	fmt.Printf("Final balance: %d\n", account.Balance())
	if account.Balance() >= 0 {
		fmt.Println("OK: No overdraft")
	} else {
		fmt.Println("BUG: Account overdrawn due to race condition!")
	}
}

// === Part 3: Workload Types ===

func ioBoundWork(id int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
	<-ch
	results <- fmt.Sprintf("I/O worker %d: done", id)
}

func cpuBoundWork(id int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	sum := 0
	for i := 0; i < 1_000_000; i++ {
		sum += i
	}
	results <- fmt.Sprintf("CPU worker %d: sum=%d", id, sum)
}

func workloadTypesDemo() {
	fmt.Println("\n--- Part 3: Workload Types (Mistake #59) ---")

	ioResults := make(chan string, 4)
	var ioWg sync.WaitGroup
	for i := 1; i <= 4; i++ {
		ioWg.Add(1)
		go ioBoundWork(i, &ioWg, ioResults)
	}
	ioWg.Wait()
	close(ioResults)

	fmt.Println("I/O-bound results (concurrent even on 1 core):")
	for msg := range ioResults {
		fmt.Println("  " + msg)
	}

	cpuResults := make(chan string, 4)
	var cpuWg sync.WaitGroup
	for i := 1; i <= 4; i++ {
		cpuWg.Add(1)
		go cpuBoundWork(i, &cpuWg, cpuResults)
	}
	cpuWg.Wait()
	close(cpuResults)

	fmt.Println("CPU-bound results (needs GOMAXPROCS > 1 for true parallelism):")
	for msg := range cpuResults {
		fmt.Println("  " + msg)
	}
}

func main() {
	fmt.Println("=== Mistakes #58-59: Race Conditions & Workload Types ===")
	fmt.Println()
	dataRaceDemo()
	raceConditionDemo()
	workloadTypesDemo()
}
