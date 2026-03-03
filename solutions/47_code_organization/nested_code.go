// nested_code.go - SOLUTION
// Refactored with early returns and guard clauses. Max nesting: 2 levels.

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

// Fixed: Guard clauses with early returns instead of deep nesting.
// Each error case is handled at the top, and the happy path stays flat.
func processOrder(order *Order) string {
	if order == nil {
		return "Error: Nil order"
	}

	if order.Customer == "" {
		return "Error: Missing customer name"
	}

	if len(order.Items) == 0 {
		return fmt.Sprintf("Order %d: No items in order", order.ID)
	}

	if order.Total <= 0 {
		return fmt.Sprintf("Order %d: Invalid total ($%.2f)", order.ID, order.Total)
	}

	if !order.Paid {
		return fmt.Sprintf("Order %d: Payment pending ($%.2f)", order.ID, order.Total)
	}

	return fmt.Sprintf("Order %d: Shipping %d items worth $%.2f to %s in %s",
		order.ID, len(order.Items), order.Total, order.Customer, order.Country)
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
