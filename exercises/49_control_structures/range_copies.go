// range_copies.go
// Understanding how range loops copy elements and when range expressions are evaluated.
//
// Mistake #30: Range loop variables are COPIES of the original elements.
//   Modifying the loop variable does NOT modify the underlying slice.
//
// Mistake #31: The range expression is evaluated ONCE before the loop starts.
//   Appending to the slice during iteration won't extend the loop.
//
// Fix both bugs so the program produces the correct output.

package main

import "fmt"

type Account struct {
	Name    string
	Balance float64
}

func main() {
	fmt.Println("=== Mistake #30: Range Copies Elements ===")

	accounts := []Account{
		{Name: "Alice", Balance: 100.0},
		{Name: "Bob", Balance: 200.0},
		{Name: "Charlie", Balance: 300.0},
	}

	// BUG: We want to add a 10% bonus to every account balance.
	// This loop modifies a COPY of each element, not the original.
	// TODO: Fix this so the original slice is actually modified.
	for _, acc := range accounts {
		acc.Balance *= 1.10
	}

	fmt.Println("After 10% bonus:")
	for _, acc := range accounts {
		fmt.Printf("  %s: $%.2f\n", acc.Name, acc.Balance)
	}
	// Expected output:
	//   Alice: $110.00
	//   Bob: $220.00
	//   Charlie: $330.00

	fmt.Println()
	fmt.Println("=== Mistake #31: Range Expression Evaluated Once ===")

	queue := []string{"task1", "task2", "task3"}
	processed := []string{}

	// BUG: This code tries to append new items to queue during iteration,
	// expecting the range loop to pick them up. It won't — the range
	// expression (len of queue) is captured once at the start.
	// TODO: Fix this to process all items including dynamically added ones.
	// Hint: Use an index-based loop with a live length check instead of range.
	for _, item := range queue {
		processed = append(processed, item)
		if item == "task2" {
			queue = append(queue, "task4")
		}
		if item == "task3" {
			queue = append(queue, "task5")
		}
	}

	fmt.Println("Processed tasks:")
	for _, p := range processed {
		fmt.Printf("  %s\n", p)
	}
	// Expected output:
	//   task1
	//   task2
	//   task3
	//   task4
	//   task5

	fmt.Println()
	fmt.Printf("Total processed: %d\n", len(processed))
	// Expected: Total processed: 5
}
