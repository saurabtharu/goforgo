// stack_heap_alloc.go
// Understand escape analysis and reduce heap allocations in Go.
//
// This exercise covers two important allocation concepts:
//
// 1. Escape analysis: The Go compiler decides whether a variable can live
//    on the stack (cheap) or must escape to the heap (expensive GC pressure).
//    When you store a pointer in a slice or map, it MUST escape to the heap
//    because the compiler can't guarantee its lifetime. Returning values
//    instead of pointers and avoiding unnecessary pointer storage keeps
//    data on the stack.
//
// 2. Reducing allocations: In hot loops, allocating on every iteration
//    creates GC pressure. Using sync.Pool to reuse objects avoids this.
//
// Fix the code to reduce heap allocations.

package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Point represents a 2D point.
type Point struct {
	X, Y float64
}

// collectPointsPtr creates points as pointers and stores them in a slice.
// Storing pointers in a slice forces every Point to escape to the heap,
// because the slice (and its contents) outlive the function that created them.
func collectPointsPtr(n int) float64 {
	points := make([]*Point, 0, n)
	for i := 0; i < n; i++ {
		p := &Point{X: float64(i), Y: float64(i + 1)}
		points = append(points, p)
	}

	total := 0.0
	for _, p := range points {
		total += p.X*p.X + p.Y*p.Y
	}
	return total
}

// collectPointsValue should create points by value, avoiding heap allocation.
// TODO: Instead of storing *Point in a slice, store Point values directly.
// A []Point slice stores values contiguously in memory - no pointer
// indirection, no per-element heap allocation, better cache locality.
func collectPointsValue(n int) float64 {
	// FIX: Change []*Point to []Point and store values instead of pointers.
	// Replace &Point{...} with Point{...} (no address-of operator).
	points := make([]*Point, 0, n)
	for i := 0; i < n; i++ {
		p := &Point{X: float64(i), Y: float64(i + 1)}
		points = append(points, p)
	}

	total := 0.0
	for _, p := range points {
		total += p.X*p.X + p.Y*p.Y
	}
	return total
}

// BufferProcessor processes byte buffers through an interface.
// Using an interface prevents the compiler from devirtualizing the call,
// which forces the buffer argument to escape to the heap.
type BufferProcessor interface {
	Process(buf []byte) byte
}

type checksummer struct{}

func (c *checksummer) Process(buf []byte) byte {
	var sum byte
	for _, b := range buf {
		sum ^= b
	}
	return sum
}

// allocateInLoop creates a new 1KB buffer on every iteration and passes
// it through an interface method. Each buffer escapes to the heap.
func allocateInLoop(proc BufferProcessor, iterations int) int {
	total := 0
	for i := 0; i < iterations; i++ {
		buf := make([]byte, 1_024)
		buf[0] = byte(i % 256)
		buf[1] = byte(i / 256 % 256)
		total += int(proc.Process(buf))
	}
	return total
}

// poolAllocateInLoop should use sync.Pool to reuse buffers.
// TODO: Create a sync.Pool with a New function that returns make([]byte, 1_024).
// In the loop, Get a buffer from the pool, use it, then Put it back.
// This dramatically reduces total bytes allocated since buffers are reused.
func poolAllocateInLoop(proc BufferProcessor, iterations int) int {
	// FIX: Create a sync.Pool:
	// pool := &sync.Pool{
	//     New: func() any { return make([]byte, 1_024) },
	// }
	_ = sync.Pool{} // placeholder - replace with real pool

	total := 0
	for i := 0; i < iterations; i++ {
		// FIX: Replace make() with pool.Get().([]byte)
		buf := make([]byte, 1_024)
		buf[0] = byte(i % 256)
		buf[1] = byte(i / 256 % 256)
		total += int(proc.Process(buf))
		// FIX: Add pool.Put(buf) here to return the buffer for reuse
	}
	return total
}

func measureAllocs(name string, fn func()) (uint64, uint64) {
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	fn()

	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	mallocs := after.Mallocs - before.Mallocs
	bytes := after.TotalAlloc - before.TotalAlloc
	fmt.Printf("  %s: %d allocs, %s total\n", name, mallocs, formatBytes(bytes))
	return mallocs, bytes
}

func formatBytes(b uint64) string {
	const (
		kb = 1_024
		mb = 1_024 * kb
	)
	switch {
	case b >= mb:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func main() {
	fmt.Println("=== Stack vs Heap Allocations ===")

	const n = 100_000

	// Part 1: Pointer slices vs value slices
	fmt.Println("\n--- Part 1: Pointer Slice vs Value Slice ---")
	fmt.Println("Storing *Point in a slice forces each Point to the heap.")
	fmt.Println("Storing Point values keeps them contiguous - fewer allocs, better cache.")

	ptrMallocs, ptrBytes := measureAllocs("Pointer slice", func() {
		_ = collectPointsPtr(n)
	})

	valMallocs, valBytes := measureAllocs("Value slice", func() {
		_ = collectPointsValue(n)
	})

	if valMallocs < ptrMallocs {
		fmt.Printf("Value slice saved %d allocations and %s!\n",
			ptrMallocs-valMallocs, formatBytes(ptrBytes-valBytes))
	} else {
		fmt.Println("HINT: Change []*Point to []Point and remove the & operator")
	}

	// Part 2: sync.Pool for buffer reuse
	fmt.Println("\n--- Part 2: Buffer Pool ---")
	fmt.Println("Allocating a buffer every iteration wastes memory.")
	fmt.Println("sync.Pool reuses buffers, reducing total bytes allocated.")

	proc := &checksummer{}
	const iters = 10_000

	_, naiveBytes := measureAllocs("Naive alloc", func() {
		_ = allocateInLoop(proc, iters)
	})

	_, poolBytes := measureAllocs("Pool alloc", func() {
		_ = poolAllocateInLoop(proc, iters)
	})

	if poolBytes < naiveBytes {
		saved := float64(naiveBytes-poolBytes) / float64(naiveBytes) * 100
		fmt.Printf("Pool reduced memory by %.0f%% (%s vs %s)!\n",
			saved, formatBytes(poolBytes), formatBytes(naiveBytes))
	} else {
		fmt.Println("HINT: Use sync.Pool to reuse buffers instead of allocating each time")
	}

	fmt.Println("\nAllocation optimization complete!")
}
