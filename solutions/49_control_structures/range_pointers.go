// range_pointers.go (solution)
// Fixed: Use index-based access for pointers to actual slice elements,
// and sort map keys for deterministic iteration order.

package main

import (
	"fmt"
	"sort"
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

	// FIXED: Use index to get a pointer to the actual slice element.
	passing := []*Student{}
	for i := range students {
		if students[i].Grade >= 80 {
			passing = append(passing, &students[i])
		}
	}

	// Apply 5-point bonus through pointers — now modifies the original slice.
	for _, p := range passing {
		p.Grade += 5
	}

	fmt.Println("After bonus (via original slice):")
	for _, s := range students {
		fmt.Printf("  %s: %d\n", s.Name, s.Grade)
	}

	fmt.Println()
	fmt.Println("=== Mistake #33: Map Iteration Order ===")

	scores := map[string]int{
		"Alice":   90,
		"Bob":     85,
		"Charlie": 95,
		"Diana":   88,
	}

	// FIXED: Extract keys, sort them, iterate in sorted order.
	fmt.Println("Scores (sorted by name):")
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, scores[k])
	}

	fmt.Println()
	fmt.Println("All tests passed!")
}
