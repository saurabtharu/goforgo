// race_conditions.go
// Understand data races vs race conditions, and I/O-bound vs CPU-bound workloads.
// Based on 100 Go Mistakes #58 and #59.
//
// Data race: concurrent unsynchronized access to shared memory (detectable by -race).
// Race condition: program correctness depends on goroutine execution order
//   (mutex alone may not fix it - you need correct LOGIC).
//
// I AM NOT DONE YET!

package main

import (
	"fmt"
	"sync"
)

// === Part 1: Data Race ===
// Multiple goroutines read and write `counter` without synchronization.
// This is a data race: undefined behavior that -race can detect.
//
// FIX: Add a sync.Mutex to protect all reads and writes to counter.

func dataRaceDemo() {
	fmt.Println("--- Part 1: Data Race ---")
	counter := 0
	// TODO: declare a sync.Mutex here

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1_000; j++ {
				// TODO: lock mutex before read-modify-write
				counter++
				// TODO: unlock mutex after modification
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

// === Part 2: Race Condition ===
// Even with a mutex, this bank transfer has a race condition.
// Two transfers check the balance independently, both see enough funds,
// and both withdraw - overdrawing the account.
//
// The bug: check-then-act is not atomic. The mutex protects individual
// operations but not the combined check+withdraw sequence.
//
// FIX: Make the entire check-and-withdraw sequence atomic by holding
// the lock for the full duration of the transfer function.

type Account struct {
	mu      sync.Mutex
	balance int
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func (a *Account) Withdraw(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance >= amount {
		a.balance -= amount
		return true
	}
	return false
}

// transfer checks balance then withdraws - but these are two separate
// critical sections! Another goroutine can withdraw between check and act.
//
// FIX: Hold the lock across the entire check-and-withdraw sequence.
// Hint: You'll need to replace the Balance()/Withdraw() calls with
// direct field access while holding the lock.
func transfer(from *Account, amount int, id int, results chan<- string) {
	// BUG: check and act are not atomic together
	if from.Balance() >= amount {
		// Another goroutine can withdraw here, between check and act!
		if from.Withdraw(amount) {
			results <- fmt.Sprintf("Transfer %d: withdrew %d, remaining: %d", id, amount, from.Balance())
		} else {
			results <- fmt.Sprintf("Transfer %d: withdraw failed (balance changed!)", id)
		}
	} else {
		results <- fmt.Sprintf("Transfer %d: insufficient funds (balance: %d)", id, from.Balance())
	}
}

func raceConditionDemo() {
	fmt.Println("\n--- Part 2: Race Condition ---")
	account := &Account{balance: 100}
	results := make(chan string, 3)

	// Three concurrent transfers of 60 each from an account with 100.
	// At most ONE should succeed, but the race condition may allow two.
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
// I/O-bound work benefits from concurrency even on a single core,
// because goroutines yield while waiting.
// CPU-bound work only benefits from true parallelism (multiple cores).
//
// This part demonstrates the difference. No fix needed - just understand
// the output and the comments.

func ioBoundWork(id int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	// Simulate I/O: channel receive blocks the goroutine but not the OS thread.
	// The Go scheduler runs other goroutines while this one waits.
	ch := make(chan struct{})
	go func() {
		// Simulate I/O completing
		ch <- struct{}{}
	}()
	<-ch
	results <- fmt.Sprintf("I/O worker %d: done", id)
}

func cpuBoundWork(id int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	// Simulate CPU work: this loop occupies the goroutine's OS thread.
	// Other goroutines on the same thread must wait.
	sum := 0
	for i := 0; i < 1_000_000; i++ {
		sum += i
	}
	results <- fmt.Sprintf("CPU worker %d: sum=%d", id, sum)
}

func workloadTypesDemo() {
	fmt.Println("\n--- Part 3: Workload Types (Mistake #59) ---")

	// I/O-bound: many goroutines can make progress concurrently
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

	// CPU-bound: goroutines compete for OS threads
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
