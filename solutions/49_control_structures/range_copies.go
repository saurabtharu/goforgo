// range_copies.go (solution)
// Fixed: Use index-based access to modify slice elements,
// and index-based loop to handle dynamically growing collections.

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

	// FIXED: Use index to modify the original slice element directly.
	for i := range accounts {
		accounts[i].Balance *= 1.10
	}

	fmt.Println("After 10% bonus:")
	for _, acc := range accounts {
		fmt.Printf("  %s: $%.2f\n", acc.Name, acc.Balance)
	}

	fmt.Println()
	fmt.Println("=== Mistake #31: Range Expression Evaluated Once ===")

	queue := []string{"task1", "task2", "task3"}
	processed := []string{}

	// FIXED: Use an index-based loop so len(queue) is re-evaluated each iteration.
	for i := 0; i < len(queue); i++ {
		item := queue[i]
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

	fmt.Println()
	fmt.Printf("Total processed: %d\n", len(processed))
}
