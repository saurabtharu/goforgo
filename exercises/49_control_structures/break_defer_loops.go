// break_defer_loops.go
// Understanding break behavior with switch/select inside loops, and defer pitfalls in loops.
//
// Mistake #34: A `break` inside a switch or select that is inside a for loop
//   only breaks out of the switch/select — NOT the enclosing loop. The loop
//   keeps running.
//
// Mistake #35: Using `defer` inside a loop accumulates deferred calls until
//   the function returns, not until the loop iteration ends. This causes
//   resource leaks when opening files/connections in a loop.
//
// Fix both bugs so the program produces the correct output.

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== Mistake #34: Break in Switch Inside Loop ===")

	// We want to process commands until we see "quit", then stop.
	commands := []string{"start", "process", "quit", "should_not_run", "also_skipped"}
	executed := []string{}

	// BUG: The `break` here only exits the switch, not the for loop.
	// The loop continues processing commands after "quit".
	// TODO: Fix this so the loop actually stops at "quit".
	for _, cmd := range commands {
		switch cmd {
		case "quit":
			fmt.Println("  Received quit command")
			break
		default:
			fmt.Printf("  Executing: %s\n", cmd)
			executed = append(executed, cmd)
		}
	}

	fmt.Printf("Executed commands: [%s]\n", strings.Join(executed, ", "))
	// Expected: Executed commands: [start, process]

	fmt.Println()
	fmt.Println("=== Mistake #35: Defer Inside Loop ===")

	// Simulate opening resources in a loop. Each "open" should be paired
	// with a "close" before the next resource is opened.
	resources := []string{"db_conn", "file_handle", "cache_client"}

	// BUG: defer inside a loop doesn't execute until the function returns.
	// All opens happen first, then all closes happen at function exit.
	// In real code this causes resource exhaustion (too many open files, etc.)
	// TODO: Fix this so each resource is closed immediately after use.
	// Hint: Extract the loop body into a separate function, or use explicit close.
	log := []string{}
	for _, name := range resources {
		log = append(log, fmt.Sprintf("open(%s)", name))
		defer func(n string) {
			log = append(log, fmt.Sprintf("close(%s)", n))
		}(name)
		log = append(log, fmt.Sprintf("use(%s)", name))
	}

	fmt.Println("Resource lifecycle:")
	for _, entry := range log {
		fmt.Printf("  %s\n", entry)
	}
	// Expected output:
	//   open(db_conn)
	//   use(db_conn)
	//   close(db_conn)
	//   open(file_handle)
	//   use(file_handle)
	//   close(file_handle)
	//   open(cache_client)
	//   use(cache_client)
	//   close(cache_client)

	fmt.Println()
	fmt.Printf("Total lifecycle events: %d\n", len(log))
	// Expected: Total lifecycle events: 9
}
