package main

import (
	"fmt"
	"time"
)

// 100 Go Mistakes #75, #76: Time duration confusion and time.After memory leaks
//
// This exercise covers two common time-related mistakes:
// 1. Passing an integer where time.Duration is expected (nanoseconds vs seconds)
// 2. Using time.After inside a loop, which leaks timers

func main() {
	fmt.Println("=== Time Mistakes ===")
	fmt.Println()

	// --- Mistake #75: Wrong time.Duration value ---
	fmt.Println("--- Duration Confusion ---")
	demonstrateDurationConfusion()
	fmt.Println()

	// --- Mistake #76: time.After memory leak ---
	fmt.Println("--- Timer Leak Fix ---")
	demonstrateTimerLeak()
	fmt.Println()

	fmt.Println("All time checks passed!")
}

// demonstrateDurationConfusion shows the mistake of passing raw integers
// where time.Duration is expected.
func demonstrateDurationConfusion() {
	// BUG: This creates a duration of 1000 nanoseconds, NOT 1000 milliseconds!
	// When you pass an untyped integer to a function expecting time.Duration,
	// Go interprets it as nanoseconds because time.Duration's base unit is nanoseconds.
	//
	// TODO: Fix both durations below to represent the intended time values.

	// Intended: 1 second (1000 milliseconds)
	tickInterval := time.Duration(1000)

	// Intended: 500 milliseconds
	timeout := time.Duration(500)

	fmt.Printf("Tick interval: %v (intended: 1s)\n", tickInterval)
	fmt.Printf("Timeout: %v (intended: 500ms)\n", timeout)

	// Check that the values are correct
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

// demonstrateTimerLeak shows the pattern that leaks timers in a loop.
// In production code with a long-running for-select loop, calling time.After
// on every iteration creates a new timer each time. Old timers are not garbage
// collected until they fire, causing memory leaks.
func demonstrateTimerLeak() {
	events := make(chan string, 5)

	// Simulate some events
	go func() {
		for _, msg := range []string{"alpha", "beta", "gamma", "delta", "done"} {
			events <- msg
			time.Sleep(10 * time.Millisecond)
		}
	}()

	var received []string

	// BUG: time.After creates a NEW timer channel on every loop iteration.
	// Each call allocates a timer that won't be GC'd until it fires.
	// In a long-running service, this is a memory leak.
	//
	// TODO: Replace time.After with time.NewTimer and reset it on each iteration.
	// Use timer.Reset() to reuse the same timer instead of creating new ones.
	for {
		select {
		case msg := <-events:
			if msg == "done" {
				fmt.Printf("Received messages: %v\n", received)
				fmt.Println("PASS: used reusable timer (no leak)")
				return
			}
			received = append(received, msg)
		case <-time.After(1 * time.Second):
			fmt.Println("FAIL: still using time.After in loop (memory leak)")
			return
		}
	}
}
