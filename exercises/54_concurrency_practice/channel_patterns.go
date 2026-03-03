// channel_patterns.go
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

// priorityReceiver should drain the high-priority channel before processing
// low-priority messages. When both channels have data, select picks randomly,
// so you can't rely on case order for priority.
//
// BUG: This implementation assumes putting highPri first in select gives it
// priority, but Go's select is random when multiple cases are ready.
//
// TODO: Fix by checking the high-priority channel in a nested select or
// by draining it first in a separate loop before checking low-priority.
func priorityReceiver(highPri <-chan string, lowPri <-chan string, results chan<- string) {
	for {
		// BUG: select chooses randomly between ready cases,
		// so high-priority messages aren't guaranteed to be processed first
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

// worker does some work and signals completion.
//
// BUG: Uses chan bool for pure signaling. The bool value is meaningless -
// it's just used as a notification that work is done.
//
// TODO: Change done channel to chan struct{} and send struct{}{} instead of true.
// chan struct{} is the idiomatic Go pattern for signaling because it uses zero memory.
func worker(id int, done chan bool, results chan<- string) {
	results <- fmt.Sprintf("worker %d: finished", id)
	done <- true // Wasteful: the bool value carries no information
}

// === Part 3: Nil Channels (Mistake #66) ===

// fanIn merges two channels into one. When one source is exhausted, it should
// stop checking that source.
//
// BUG: When a channel is closed, the receive returns the zero value immediately.
// Without setting closed channels to nil, the select keeps hitting the closed
// channel case and producing zero-value results.
//
// TODO: Set a channel to nil after it's closed. A nil channel in select
// blocks forever, effectively disabling that case.
func fanIn(ch1, ch2 <-chan int) []int {
	var results []int
	active := 2

	for active > 0 {
		select {
		case v, ok := <-ch1:
			if !ok {
				active--
				// BUG: ch1 is closed but not set to nil.
				// select will keep picking this case, returning zero values.
				continue
			}
			results = append(results, v)
		case v, ok := <-ch2:
			if !ok {
				active--
				// BUG: same problem as ch1
				continue
			}
			results = append(results, v)
		}
	}
	return results
}

// === Part 4: Channel Sizing (Mistake #67) ===

// produceValues creates goroutines that send values on a channel.
//
// BUG: The channel is unbuffered, so goroutines block on send until someone
// reads. If there are more producers than the consumer reads, goroutines leak.
//
// TODO: Make the channel buffered with capacity equal to the number of producers.
func produceValues(n int) []int {
	// BUG: Unbuffered channel - producers block until consumer reads each value
	ch := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			ch <- val // Blocks on unbuffered channel
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

	// BUG: Using chan bool when chan struct{} is more appropriate
	done := make(chan bool, 3)
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
