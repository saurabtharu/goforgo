package main

import "fmt"

// This exercise demonstrates Mistake #47: defer arguments are evaluated
// immediately when the defer statement is executed, not when the deferred
// function actually runs.

func logStatus(step string, status string) {
	fmt.Printf("Step: %s, Status: %s\n", step, status)
}

// Fixed: Use a closure to capture status by reference.
// The closure doesn't evaluate status until it runs at function exit.
func processSteps() {
	status := "pending"

	// Fixed: closure captures status by reference, seeing the final value
	defer func() { logStatus("cleanup", status) }()

	fmt.Println("Running step 1...")
	status = "step1_done"

	fmt.Println("Running step 2...")
	status = "step2_done"

	fmt.Println("Running step 3...")
	status = "completed"

	fmt.Printf("Final status before return: %s\n", status)
}

// Fixed: Use a closure to capture counter by reference.
func trackCounter() {
	counter := 0

	// Fixed: closure captures counter by reference
	defer func() { fmt.Printf("Final counter: %d\n", counter) }()

	for i := 0; i < 5; i++ {
		counter++
		fmt.Printf("Counter: %d\n", counter)
	}
}

// Fixed: Use a closure for the sum defer to capture the final value.
func deferInLoop() {
	sum := 0

	// Fixed: closure captures sum by reference, prints the final sum
	defer func() { fmt.Printf("Total sum: %d\n", sum) }()

	for i := 1; i <= 3; i++ {
		sum += i
		// This correctly captures i's current value at each iteration.
		// Here, immediate evaluation is actually what we want.
		defer fmt.Printf("Deferred value: %d\n", i)
	}
}

func main() {
	fmt.Println("=== processSteps ===")
	processSteps()

	fmt.Println()
	fmt.Println("=== trackCounter ===")
	trackCounter()

	fmt.Println()
	fmt.Println("=== deferInLoop ===")
	deferInLoop()
}
