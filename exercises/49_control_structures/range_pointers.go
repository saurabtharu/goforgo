// range_pointers.go
// Understanding pointer semantics with range loops and map iteration behavior.
//
// Mistake #32: Range variables are copies. Even in Go 1.22+ (where each
//   iteration gets its own variable), taking &v gives you a pointer to a
//   COPY, not to the original slice element. Modifying through that pointer
//   does NOT modify the original slice.
//
// Mistake #33: Map iteration order is non-deterministic in Go. Code that
//   depends on a specific traversal order will produce inconsistent output.
//
// Fix both bugs so the program produces the correct output.

package main

import (
	"fmt"
	// TODO: You'll need the "sort" package to fix Mistake #33
)

type Student struct {
	Name  string
	Grade int
}

func main() {
	fmt.Println("=== Mistake #32: Pointers to Range Copies ===")

	students := []Student{
		{Name: "Alice", Grade: 85},
		{Name: "Bob", Grade: 72},
		{Name: "Charlie", Grade: 91},
	}

	// We want to collect pointers to students who pass (grade >= 80)
	// AND bump their grades by 5 bonus points through those pointers.
	// BUG: &s points to a copy of the element, not the original.
	// Modifying *ptr changes the copy, not students[i].
	// TODO: Fix this so pointers refer to actual slice elements and
	// modifications are reflected in the original slice.
	passing := []*Student{}
	for _, s := range students {
		if s.Grade >= 80 {
			passing = append(passing, &s)
		}
	}

	// Apply 5-point bonus through pointers
	for _, p := range passing {
		p.Grade += 5
	}

	fmt.Println("After bonus (via original slice):")
	for _, s := range students {
		fmt.Printf("  %s: %d\n", s.Name, s.Grade)
	}
	// Expected output (bonus applied to Alice and Charlie):
	//   Alice: 90
	//   Bob: 72
	//   Charlie: 96

	fmt.Println()
	fmt.Println("=== Mistake #33: Map Iteration Order ===")

	scores := map[string]int{
		"Alice":   90,
		"Bob":     85,
		"Charlie": 95,
		"Diana":   88,
	}

	// BUG: This code prints map entries expecting alphabetical order,
	// but map iteration in Go is deliberately randomized.
	// TODO: Fix this to produce deterministic alphabetical output.
	// Hint: Extract the keys, sort them, and iterate in sorted order.
	fmt.Println("Scores (sorted by name):")
	for name, score := range scores {
		fmt.Printf("  %s: %d\n", name, score)
	}
	// Expected output (always in this order):
	//   Alice: 90
	//   Bob: 85
	//   Charlie: 95
	//   Diana: 88

	fmt.Println()
	fmt.Println("All tests passed!")
}
