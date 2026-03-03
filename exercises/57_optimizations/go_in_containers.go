// go_in_containers.go
// Understand GOMAXPROCS behavior in containerized environments.
//
// When a Go program runs inside a container (Docker, Kubernetes), the
// runtime.NumCPU() function returns the HOST machine's CPU count, not the
// container's CPU limit. This means GOMAXPROCS defaults to the host CPU count,
// which can cause excessive goroutine scheduling overhead when the container
// is limited to fewer cores.
//
// For example, a container limited to 2 CPUs on a 64-core host will set
// GOMAXPROCS=64, creating 64 OS threads competing for 2 actual CPU cores.
//
// Fix this program to properly detect and set GOMAXPROCS based on container
// CPU limits, or a reasonable fallback.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// detectContainerCPULimit attempts to read the CPU limit from the
// container's cgroup filesystem. In cgroups v2 (modern Linux/Docker),
// the limit is in /sys/fs/cgroup/cpu.max.
//
// Returns the detected limit, or 0 if not in a container or unreadable.
func detectContainerCPULimit() int {
	// Try cgroups v2 (modern containers)
	data, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err == nil {
		parts := strings.Fields(string(data))
		if len(parts) >= 2 && parts[0] != "max" {
			quota, err1 := strconv.Atoi(parts[0])
			period, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil && period > 0 {
				cpus := quota / period
				if cpus < 1 {
					cpus = 1
				}
				return cpus
			}
		}
	}

	// Try cgroups v1 (older containers)
	quotaData, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err == nil {
		periodData, err := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
		if err == nil {
			quota, err1 := strconv.Atoi(strings.TrimSpace(string(quotaData)))
			period, err2 := strconv.Atoi(strings.TrimSpace(string(periodData)))
			if err1 == nil && err2 == nil && quota > 0 && period > 0 {
				cpus := quota / period
				if cpus < 1 {
					cpus = 1
				}
				return cpus
			}
		}
	}

	return 0
}

// setContainerAwareGOMAXPROCS sets GOMAXPROCS appropriately.
// TODO: This function currently does nothing. Fix it to:
// 1. Check if there's an explicit GOMAXPROCS environment variable set
// 2. Try to detect the container CPU limit using detectContainerCPULimit()
// 3. Fall back to runtime.NumCPU() if not in a container
// 4. Actually call runtime.GOMAXPROCS() with the determined value
func setContainerAwareGOMAXPROCS() int {
	// FIX: Check for explicit GOMAXPROCS environment variable first.
	// If the user set GOMAXPROCS explicitly, respect that.

	// FIX: Try container detection. If detectContainerCPULimit() returns
	// a positive number, use that value.

	// FIX: Fall back to runtime.NumCPU() if we're not in a container.

	// FIX: Call runtime.GOMAXPROCS(cpus) with the determined value
	// and return the value that was set.

	return runtime.NumCPU() // This is the WRONG default for containers
}

// simulateWorkload runs a parallel workload and measures how it performs
// with the current GOMAXPROCS setting.
func simulateWorkload(workers int) time.Duration {
	var wg sync.WaitGroup
	wg.Add(workers)

	start := time.Now()
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// Simulate CPU-bound work
			sum := 0
			for j := 0; j < 1_000_000; j++ {
				sum += j
			}
			_ = sum
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	fmt.Println("=== GOMAXPROCS in Containers ===")

	// Show the problem: runtime.NumCPU() returns host CPU count
	fmt.Println("\n--- Host vs Container CPU Detection ---")
	fmt.Printf("runtime.NumCPU() reports: %d CPUs\n", runtime.NumCPU())
	fmt.Printf("Current GOMAXPROCS:       %d\n", runtime.GOMAXPROCS(0))

	// Detect container limits
	containerLimit := detectContainerCPULimit()
	if containerLimit > 0 {
		fmt.Printf("Container CPU limit:      %d CPUs\n", containerLimit)
	} else {
		fmt.Println("Container CPU limit:      not detected (not in container or no limit)")
	}

	// Check for environment variable
	envGOMAXPROCS := os.Getenv("GOMAXPROCS")
	if envGOMAXPROCS != "" {
		fmt.Printf("GOMAXPROCS env var:       %s\n", envGOMAXPROCS)
	} else {
		fmt.Println("GOMAXPROCS env var:       not set")
	}

	// Set container-aware GOMAXPROCS
	fmt.Println("\n--- Setting Container-Aware GOMAXPROCS ---")
	effectiveCPUs := setContainerAwareGOMAXPROCS()
	fmt.Printf("Effective GOMAXPROCS:     %d\n", effectiveCPUs)
	fmt.Printf("Actual GOMAXPROCS:        %d\n", runtime.GOMAXPROCS(0))

	// Demonstrate workload with different GOMAXPROCS values
	fmt.Println("\n--- Workload Performance ---")

	// Test with GOMAXPROCS=1
	runtime.GOMAXPROCS(1)
	time1 := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=1:  4 workers took %v\n", time1)

	// Test with GOMAXPROCS=2
	runtime.GOMAXPROCS(2)
	time2 := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=2:  4 workers took %v\n", time2)

	// Test with effective GOMAXPROCS
	runtime.GOMAXPROCS(effectiveCPUs)
	timeN := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=%d: 4 workers took %v\n", effectiveCPUs, timeN)

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Goroutines active: %d\n", runtime.NumGoroutine())
	fmt.Println("In containers, always set GOMAXPROCS to match the container's CPU limit.")
	fmt.Println("Use automaxprocs (uber-go) in production, or set GOMAXPROCS env var.")

	fmt.Println("\nContainer optimization complete!")
}
