// nested_code.go
// Refactor deeply nested code using early returns and guard clauses
//
// Deeply nested code is hard to read and maintain. The "happy path" gets
// buried inside layers of indentation. Go's convention is to handle errors
// and edge cases early with guard clauses, keeping the main logic at the
// lowest indentation level.
//
// Refactor the processOrder function to use early returns instead of deep
// nesting. The output must remain exactly the same.

package main

import "fmt"

type Order struct {
	ID       int
	Customer string
	Items    []string
	Total    float64
	Paid     bool
	Country  string
}

// TODO: Refactor this function to use early returns and guard clauses.
// The function currently has 5 levels of nesting. Flatten it so the
// deepest nesting is at most 2 levels (function body + one if/for).
// The output must remain exactly the same!
func processOrder(order *Order) string {
	if order != nil {
		if order.Customer != "" {
			if len(order.Items) > 0 {
				if order.Total > 0 {
					if order.Paid {
						result := fmt.Sprintf("Order %d: Shipping %d items worth $%.2f to %s in %s",
							order.ID, len(order.Items), order.Total, order.Customer, order.Country)
						return result
					} else {
						return fmt.Sprintf("Order %d: Payment pending ($%.2f)", order.ID, order.Total)
					}
				} else {
					return fmt.Sprintf("Order %d: Invalid total ($%.2f)", order.ID, order.Total)
				}
			} else {
				return fmt.Sprintf("Order %d: No items in order", order.ID)
			}
		} else {
			return "Error: Missing customer name"
		}
	} else {
		return "Error: Nil order"
	}
}

func main() {
	fmt.Println("=== Nested Code Refactoring ===")

	orders := []*Order{
		nil,
		{ID: 1, Customer: "", Items: []string{"book"}, Total: 15.99, Paid: true, Country: "US"},
		{ID: 2, Customer: "Alice", Items: []string{}, Total: 25.00, Paid: true, Country: "UK"},
		{ID: 3, Customer: "Bob", Items: []string{"pen", "notebook"}, Total: -5.00, Paid: true, Country: "CA"},
		{ID: 4, Customer: "Charlie", Items: []string{"laptop"}, Total: 999.99, Paid: false, Country: "DE"},
		{ID: 5, Customer: "Diana", Items: []string{"phone", "case", "charger"}, Total: 549.50, Paid: true, Country: "JP"},
	}

	for _, order := range orders {
		result := processOrder(order)
		fmt.Println(result)
	}

	fmt.Println("\nAll orders processed!")
}
