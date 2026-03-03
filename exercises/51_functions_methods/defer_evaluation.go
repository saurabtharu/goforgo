package main

import "fmt"

// This exercise demonstrates Mistake #47: defer arguments are evaluated
// immediately when the defer statement is executed, not when the deferred
// function actually runs.

// logStatus prints the status of an operation.
func logStatus(step string, status string) {
	fmt.Printf("Step: %s, Status: %s\n", step, status)
}

// processSteps runs a series of steps and defers a status log.
// BUG: The defer captures 'status' by value at defer time, which is "pending".
// When the deferred function runs at function exit, it prints "pending"
// instead of the final value of status.
// FIX: Use a closure to capture status by reference so the deferred function
// sees the final value when it executes.
func processSteps() {
	status := "pending"

	// BUG: defer evaluates logStatus's arguments NOW, capturing status="pending"
	defer logStatus("cleanup", status)

	fmt.Println("Running step 1...")
	status = "step1_done"

	fmt.Println("Running step 2...")
	status = "step2_done"

	fmt.Println("Running step 3...")
	status = "completed"

	fmt.Printf("Final status before return: %s\n", status)
}

// trackCounter demonstrates the same trap with integer arguments.
// BUG: defer captures counter's value (0) at defer time, not at function exit.
// FIX: Use a closure defer func() { ... }() to capture by reference.
func trackCounter() {
	counter := 0

	// BUG: This captures counter=0, not the final value
	defer fmt.Printf("Final counter: %d\n", counter)

	for i := 0; i < 5; i++ {
		counter++
		fmt.Printf("Counter: %d\n", counter)
	}
}

// deferInLoop demonstrates correct defer usage with closures in a loop.
// BUG: The defer captures 'i' by value at each iteration, so this actually
// works correctly for printing each value. But the "sum" defer has the
// same evaluation trap.
// FIX: Use a closure for the sum defer to capture the final value.
func deferInLoop() {
	sum := 0

	// BUG: This captures sum=0, printed last due to defer LIFO order
	defer fmt.Printf("Total sum: %d\n", sum)

	for i := 1; i <= 3; i++ {
		sum += i
		// This correctly captures i's current value (defer arg evaluation is useful here)
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
