// channels_vs_mutexes.go
// Learn when to use channels vs mutexes - the right tool for the right job.
// Based on 100 Go Mistakes #57.
//
// Rule of thumb:
//   - Mutex: protecting shared state (counters, caches, maps)
//   - Channels: communicating between goroutines (pipelines, signals)
//
// This exercise has two parts that use the WRONG synchronization primitive.
// Your job: rewrite each part with the appropriate one.
//
// I AM NOT DONE YET!

package main

import (
	"fmt"
	"sync"
)

// === Part 1: Shared Counter ===
// This counter uses channels to coordinate increments.
// It works, but it's overcomplicated for simple shared state.
//
// TODO: Replace this entire channel-based approach with a struct that uses
// sync.Mutex to protect the counter. The struct should have:
//   - mu sync.Mutex
//   - value int
//   - An Increment() method that locks, increments, and unlocks
//   - A Value() method that locks, reads, and unlocks

func counterWithChannels(numWorkers, incrementsPerWorker int) int {
	type increment struct{}
	ch := make(chan increment, 100)
	resultCh := make(chan int)

	// Dedicated goroutine to process all increments sequentially.
	// This is a lot of machinery for a simple counter!
	go func() {
		count := 0
		for range ch {
			count++
		}
		resultCh <- count
	}()

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerWorker; j++ {
				ch <- increment{}
			}
		}()
	}

	wg.Wait()
	close(ch)
	return <-resultCh
}

// === Part 2: Data Pipeline ===
// This pipeline uses mutexes and shared slices to pass data between stages.
// It works, but it's unnatural - channels are the right tool for pipelines.
//
// TODO: Replace this mutex-based approach with a channel pipeline:
//   - generate() returns a <-chan int (sends values, then closes)
//   - double(in <-chan int) returns a <-chan int (reads, doubles, sends, closes)
//   - collect(in <-chan int) returns []int (drains the channel into a slice)
// The pipeline becomes: collect(double(generate(values)))

func pipelineWithMutex(input []int) (doubled []int, sum int) {
	var mu sync.Mutex
	var stage1Done, stage2Done bool
	shared := make([]int, 0, len(input))
	output := make([]int, 0, len(input))

	var wg sync.WaitGroup

	// Stage 1: load input into shared slice
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		for _, v := range input {
			shared = append(shared, v)
		}
		stage1Done = true
		mu.Unlock()
	}()

	// Stage 2: read shared slice, double values, write to output
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Spin-wait until stage 1 is done (ugly!)
		for {
			mu.Lock()
			if stage1Done {
				for _, v := range shared {
					output = append(output, v*2)
				}
				stage2Done = true
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}()

	// Stage 3: read output slice, compute sum
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			mu.Lock()
			if stage2Done {
				for _, v := range output {
					sum += v
				}
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}()

	wg.Wait()
	doubled = output
	return
}

func main() {
	fmt.Println("=== Mistake #57: Channels vs Mutexes ===")
	fmt.Println()

	// Part 1: Counter (should use mutex, currently uses channels)
	fmt.Println("--- Part 1: Shared State (Counter) ---")
	result := counterWithChannels(5, 1_000)
	fmt.Printf("Counter result: %d (expected: 5000)\n", result)
	if result == 5_000 {
		fmt.Println("OK: Counter is correct")
	}
	fmt.Println("NOTE: Channels work but are overcomplicated for shared state.")
	fmt.Println("TODO: Rewrite using sync.Mutex")
	fmt.Println()

	// Part 2: Pipeline (should use channels, currently uses mutexes)
	fmt.Println("--- Part 2: Communication (Pipeline) ---")
	input := []int{1, 2, 3, 4, 5}
	doubled, sum := pipelineWithMutex(input)
	fmt.Printf("Input:   %v\n", input)
	fmt.Printf("Doubled: %v\n", doubled)
	fmt.Printf("Sum:     %d (expected: 30)\n", sum)
	if sum == 30 {
		fmt.Println("OK: Pipeline is correct")
	}
	fmt.Println("NOTE: Mutexes work but spin-waiting is wasteful for pipelines.")
	fmt.Println("TODO: Rewrite using channels")
}
