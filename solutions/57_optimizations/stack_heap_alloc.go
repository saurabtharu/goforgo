// stack_heap_alloc.go - SOLUTION
// Heap allocations reduced through value slices and sync.Pool.

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

// collectPointsPtr stores pointers (forces each Point to heap).
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

// collectPointsValue stores values directly - no per-element heap allocs.
// Fixed: []Point instead of []*Point, and Point{} instead of &Point{}.
func collectPointsValue(n int) float64 {
	points := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		p := Point{X: float64(i), Y: float64(i + 1)}
		points = append(points, p)
	}

	total := 0.0
	for _, p := range points {
		total += p.X*p.X + p.Y*p.Y
	}
	return total
}

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

// poolAllocateInLoop reuses buffers via sync.Pool.
// Fixed: pool.Get() reuses existing buffers, pool.Put() returns them.
func poolAllocateInLoop(proc BufferProcessor, iterations int) int {
	pool := &sync.Pool{
		New: func() any {
			return make([]byte, 1_024)
		},
	}

	total := 0
	for i := 0; i < iterations; i++ {
		buf := pool.Get().([]byte)
		buf[0] = byte(i % 256)
		buf[1] = byte(i / 256 % 256)
		total += int(proc.Process(buf))
		pool.Put(buf)
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
