// channel_patterns.go - SOLUTION
// Learn correct channel usage patterns and avoid common channel mistakes
//
// Mistakes #64-67 from "100 Go Mistakes":
// #64: Non-deterministic select behavior (expecting case order to matter)
// #65: Not using notification channels (chan struct{} for signaling)
// #66: Not using nil channels to disable select cases
// #67: Channel sizing mistakes causing goroutine leaks

package main

import (
	"fmt"
	"sync"
)

// === Part 1: Select Priority (Mistake #64) ===

// priorityReceiver drains the high-priority channel first before considering
// low-priority messages. Uses a nested select pattern to enforce priority.
func priorityReceiver(highPri <-chan string, lowPri <-chan string, results chan<- string) {
	for {
		// FIX: Always try high priority first in a non-blocking check
		select {
		case msg, ok := <-highPri:
			if !ok {
				highPri = nil
			} else {
				results <- fmt.Sprintf("[HIGH] %s", msg)
				continue
			}
		default:
			// No high-priority message ready, fall through
		}

		// Only check low priority if high priority is empty
		if highPri == nil && lowPri == nil {
			return
		}

		select {
		case msg, ok := <-highPri:
			if !ok {
				highPri = nil
			} else {
				results <- fmt.Sprintf("[HIGH] %s", msg)
			}
		case msg, ok := <-lowPri:
			if !ok {
				lowPri = nil
			} else {
				results <- fmt.Sprintf("[LOW] %s", msg)
			}
		}

		if highPri == nil && lowPri == nil {
			return
		}
	}
}

// === Part 2: Notification Channels (Mistake #65) ===

// worker does some work and signals completion using chan struct{} (zero-size).
func worker(id int, done chan struct{}, results chan<- string) {
	results <- fmt.Sprintf("worker %d: finished", id)
	done <- struct{}{} // FIX: Use struct{}{} for zero-cost signaling
}

// === Part 3: Nil Channels (Mistake #66) ===

// fanIn merges two channels into one. When a channel is closed, it's set to nil
// to disable that select case.
func fanIn(ch1, ch2 <-chan int) []int {
	var results []int
	active := 2

	for active > 0 {
		select {
		case v, ok := <-ch1:
			if !ok {
				active--
				ch1 = nil // FIX: Set to nil so select ignores this case
				continue
			}
			results = append(results, v)
		case v, ok := <-ch2:
			if !ok {
				active--
				ch2 = nil // FIX: Set to nil so select ignores this case
				continue
			}
			results = append(results, v)
		}
	}
	return results
}

// === Part 4: Channel Sizing (Mistake #67) ===

// produceValues creates goroutines that send values on a buffered channel.
func produceValues(n int) []int {
	// FIX: Buffer the channel to match the number of producers
	// This prevents goroutines from blocking on send
	ch := make(chan int, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			ch <- val
		}(i)
	}

	// Close channel after all producers finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []int
	for v := range ch {
		results = append(results, v)
	}
	return results
}

func main() {
	fmt.Println("=== Select Priority (Mistake #64) ===")

	highPri := make(chan string, 5)
	lowPri := make(chan string, 5)
	results := make(chan string, 10)

	// Load both channels before starting receiver
	highPri <- "URGENT-1"
	highPri <- "URGENT-2"
	highPri <- "URGENT-3"
	close(highPri)

	lowPri <- "normal-1"
	lowPri <- "normal-2"
	close(lowPri)

	go func() {
		priorityReceiver(highPri, lowPri, results)
		close(results)
	}()

	fmt.Println("Messages received (high priority should come first):")
	for msg := range results {
		fmt.Println(msg)
	}

	fmt.Println()
	fmt.Println("=== Notification Channels (Mistake #65) ===")

	// FIX: Use chan struct{} for signaling - zero memory cost
	done := make(chan struct{}, 3)
	notifyResults := make(chan string, 3)

	for i := 1; i <= 3; i++ {
		go worker(i, done, notifyResults)
	}

	// Wait for all workers
	for i := 0; i < 3; i++ {
		<-done
	}
	close(notifyResults)

	for msg := range notifyResults {
		fmt.Println(msg)
	}
	fmt.Println("all workers signaled completion")

	fmt.Println()
	fmt.Println("=== Nil Channels (Mistake #66) ===")

	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 10
	ch1 <- 20
	close(ch1)

	ch2 <- 30
	ch2 <- 40
	ch2 <- 50
	close(ch2)

	merged := fanIn(ch1, ch2)
	fmt.Printf("merged values: %v\n", merged)
	fmt.Printf("expected 5 values, got %d\n", len(merged))
	if len(merged) == 5 {
		fmt.Println("PASS: Correct number of values merged")
	} else {
		fmt.Println("FAIL: Got wrong number of values (nil channel issue)")
	}

	fmt.Println()
	fmt.Println("=== Channel Sizing (Mistake #67) ===")

	produced := produceValues(5)
	fmt.Printf("produced: %v\n", produced)
	if len(produced) == 5 {
		fmt.Println("PASS: All values received without goroutine leaks")
	} else {
		fmt.Printf("FAIL: Expected 5, got %d\n", len(produced))
	}
}
