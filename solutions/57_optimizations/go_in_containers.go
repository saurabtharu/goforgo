// go_in_containers.go - SOLUTION
// Proper container-aware GOMAXPROCS setting implemented.

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

func detectContainerCPULimit() int {
	// Try cgroups v2
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

	// Try cgroups v1
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

// setContainerAwareGOMAXPROCS sets GOMAXPROCS with proper container awareness.
// Fixed: checks env var first, then container limit, then falls back to NumCPU.
func setContainerAwareGOMAXPROCS() int {
	// 1. Respect explicit GOMAXPROCS environment variable
	if envVal := os.Getenv("GOMAXPROCS"); envVal != "" {
		if n, err := strconv.Atoi(envVal); err == nil && n > 0 {
			runtime.GOMAXPROCS(n)
			return n
		}
	}

	// 2. Try container CPU limit detection
	if limit := detectContainerCPULimit(); limit > 0 {
		runtime.GOMAXPROCS(limit)
		return limit
	}

	// 3. Fall back to host CPU count (appropriate for non-container environments)
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	return cpus
}

func simulateWorkload(workers int) time.Duration {
	var wg sync.WaitGroup
	wg.Add(workers)

	start := time.Now()
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
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

	fmt.Println("\n--- Host vs Container CPU Detection ---")
	fmt.Printf("runtime.NumCPU() reports: %d CPUs\n", runtime.NumCPU())
	fmt.Printf("Current GOMAXPROCS:       %d\n", runtime.GOMAXPROCS(0))

	containerLimit := detectContainerCPULimit()
	if containerLimit > 0 {
		fmt.Printf("Container CPU limit:      %d CPUs\n", containerLimit)
	} else {
		fmt.Println("Container CPU limit:      not detected (not in container or no limit)")
	}

	envGOMAXPROCS := os.Getenv("GOMAXPROCS")
	if envGOMAXPROCS != "" {
		fmt.Printf("GOMAXPROCS env var:       %s\n", envGOMAXPROCS)
	} else {
		fmt.Println("GOMAXPROCS env var:       not set")
	}

	fmt.Println("\n--- Setting Container-Aware GOMAXPROCS ---")
	effectiveCPUs := setContainerAwareGOMAXPROCS()
	fmt.Printf("Effective GOMAXPROCS:     %d\n", effectiveCPUs)
	fmt.Printf("Actual GOMAXPROCS:        %d\n", runtime.GOMAXPROCS(0))

	fmt.Println("\n--- Workload Performance ---")

	runtime.GOMAXPROCS(1)
	time1 := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=1:  4 workers took %v\n", time1)

	runtime.GOMAXPROCS(2)
	time2 := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=2:  4 workers took %v\n", time2)

	runtime.GOMAXPROCS(effectiveCPUs)
	timeN := simulateWorkload(4)
	fmt.Printf("GOMAXPROCS=%d: 4 workers took %v\n", effectiveCPUs, timeN)

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Goroutines active: %d\n", runtime.NumGoroutine())
	fmt.Println("In containers, always set GOMAXPROCS to match the container's CPU limit.")
	fmt.Println("Use automaxprocs (uber-go) in production, or set GOMAXPROCS env var.")

	fmt.Println("\nContainer optimization complete!")
}
