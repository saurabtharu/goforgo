// break_defer_loops.go (solution)
// Fixed: Use labeled break to exit loop from switch,
// and extract loop body to a function so defer runs per iteration.

package main

import (
	"fmt"
	"strings"
)

// useResource handles a single resource lifecycle, so defer runs
// at the end of each call instead of accumulating.
// Named return lets defer modify the result before it's returned.
func useResource(name string) (log []string) {
	log = append(log, fmt.Sprintf("open(%s)", name))
	// In real code you'd defer the actual close here.
	// We simulate it by appending to the log.
	defer func() {
		log = append(log, fmt.Sprintf("close(%s)", name))
	}()
	log = append(log, fmt.Sprintf("use(%s)", name))
	return
}

func main() {
	fmt.Println("=== Mistake #34: Break in Switch Inside Loop ===")

	commands := []string{"start", "process", "quit", "should_not_run", "also_skipped"}
	executed := []string{}

	// FIXED: Use a labeled loop so break exits the for loop, not just the switch.
loop:
	for _, cmd := range commands {
		switch cmd {
		case "quit":
			fmt.Println("  Received quit command")
			break loop
		default:
			fmt.Printf("  Executing: %s\n", cmd)
			executed = append(executed, cmd)
		}
	}

	fmt.Printf("Executed commands: [%s]\n", strings.Join(executed, ", "))

	fmt.Println()
	fmt.Println("=== Mistake #35: Defer Inside Loop ===")

	resources := []string{"db_conn", "file_handle", "cache_client"}

	// FIXED: Extract loop body into useResource() so defer runs per call.
	log := []string{}
	for _, name := range resources {
		entries := useResource(name)
		log = append(log, entries...)
	}

	fmt.Println("Resource lifecycle:")
	for _, entry := range log {
		fmt.Printf("  %s\n", entry)
	}

	fmt.Println()
	fmt.Printf("Total lifecycle events: %d\n", len(log))
}
