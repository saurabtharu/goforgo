// channels_vs_mutexes.go - SOLUTION
// Learn when to use channels vs mutexes - the right tool for the right job.
// Based on 100 Go Mistakes #57.

package main

import (
	"fmt"
	"sync"
)

// === Part 1: Shared Counter with Mutex ===
// Mutex is the right tool for protecting shared state.
// Simple, clear, no extra goroutines needed.

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func counterWithMutex(numWorkers, incrementsPerWorker int) int {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerWorker; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()
	return counter.Value()
}

// === Part 2: Data Pipeline with Channels ===
// Channels are the right tool for communication between goroutines.
// Each stage is a goroutine connected by channels - no spin-waiting.

func generate(values []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range values {
			out <- v
		}
		close(out)
	}()
	return out
}

func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * 2
		}
		close(out)
	}()
	return out
}

func collect(in <-chan int) []int {
	var result []int
	for v := range in {
		result = append(result, v)
	}
	return result
}

func pipelineWithChannels(input []int) (doubled []int, sum int) {
	// Clean pipeline: generate -> double -> collect
	doubled = collect(double(generate(input)))
	for _, v := range doubled {
		sum += v
	}
	return
}

func main() {
	fmt.Println("=== Mistake #57: Channels vs Mutexes ===")
	fmt.Println()

	// Part 1: Counter with mutex (correct tool for shared state)
	fmt.Println("--- Part 1: Shared State (Counter) ---")
	result := counterWithMutex(5, 1_000)
	fmt.Printf("Counter result: %d (expected: 5000)\n", result)
	if result == 5_000 {
		fmt.Println("OK: Counter is correct")
	}
	fmt.Println("Using sync.Mutex: simple, clear, no extra goroutines needed.")
	fmt.Println()

	// Part 2: Pipeline with channels (correct tool for communication)
	fmt.Println("--- Part 2: Communication (Pipeline) ---")
	input := []int{1, 2, 3, 4, 5}
	doubled, sum := pipelineWithChannels(input)
	fmt.Printf("Input:   %v\n", input)
	fmt.Printf("Doubled: %v\n", doubled)
	fmt.Printf("Sum:     %d (expected: 30)\n", sum)
	if sum == 30 {
		fmt.Println("OK: Pipeline is correct")
	}
	fmt.Println("Using channels: natural flow, no spin-waiting, composable stages.")
}
