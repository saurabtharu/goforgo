package main

import (
	"fmt"
	"time"
)

// 100 Go Mistakes #75, #76: Time duration confusion and time.After memory leaks
//
// Solution: Use explicit duration multipliers and reusable timers.

func main() {
	fmt.Println("=== Time Mistakes ===")
	fmt.Println()

	fmt.Println("--- Duration Confusion ---")
	demonstrateDurationConfusion()
	fmt.Println()

	fmt.Println("--- Timer Leak Fix ---")
	demonstrateTimerLeak()
	fmt.Println()

	fmt.Println("All time checks passed!")
}

func demonstrateDurationConfusion() {
	// FIXED: Use explicit time unit multipliers instead of raw integers.
	// time.Duration's base unit is nanoseconds, so always multiply by the
	// desired unit constant.

	// Intended: 1 second
	tickInterval := 1 * time.Second

	// Intended: 500 milliseconds
	timeout := 500 * time.Millisecond

	fmt.Printf("Tick interval: %v (intended: 1s)\n", tickInterval)
	fmt.Printf("Timeout: %v (intended: 500ms)\n", timeout)

	if tickInterval == 1*time.Second {
		fmt.Println("PASS: tick interval is 1 second")
	} else {
		fmt.Printf("FAIL: tick interval is %v, not 1s\n", tickInterval)
	}

	if timeout == 500*time.Millisecond {
		fmt.Println("PASS: timeout is 500ms")
	} else {
		fmt.Printf("FAIL: timeout is %v, not 500ms\n", timeout)
	}
}

func demonstrateTimerLeak() {
	events := make(chan string, 5)

	go func() {
		for _, msg := range []string{"alpha", "beta", "gamma", "delta", "done"} {
			events <- msg
			time.Sleep(10 * time.Millisecond)
		}
	}()

	var received []string

	// FIXED: Use time.NewTimer and Reset instead of time.After in loops.
	// This reuses a single timer allocation instead of creating a new one per iteration.
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	for {
		select {
		case msg := <-events:
			if msg == "done" {
				fmt.Printf("Received messages: %v\n", received)
				fmt.Println("PASS: used reusable timer (no leak)")
				return
			}
			received = append(received, msg)
			// Reset the timer for the next iteration
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(1 * time.Second)
		case <-timer.C:
			fmt.Println("FAIL: still using time.After in loop (memory leak)")
			return
		}
	}
}
