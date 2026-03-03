// generics_and_any.go - SOLUTION
// Replaced all any/interface{} usage with proper generic type constraints.

package main

import (
	"fmt"
	"strings"
)

// Fixed: Type constraint for numeric types.
type Number interface {
	int | float64
}

// Fixed: Generic sum with type constraint. Compile-time type safety.
func sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// Fixed: Generic contains with comparable constraint.
func contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// Fixed: Type constraint for ordered types (supports <).
type Ordered interface {
	~int | ~float64 | ~string
}

// Fixed: Generic minimum with Ordered constraint. No runtime type assertions.
func minimum[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Fixed: Generic map with two type parameters.
func mapSlice[T any, U any](input []T, fn func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func main() {
	fmt.Println("=== Generics and Any ===")

	// Test 1: Sum - now type-safe
	fmt.Println("\n--- Sum ---")
	ints := []int{1, 2, 3, 4, 5}
	fmt.Println("Sum of ints:", sum(ints))

	floats := []float64{1.5, 2.5, 3.5}
	fmt.Println("Sum of floats:", sum(floats))

	// Test 2: Contains - now type-safe
	fmt.Println("\n--- Contains ---")
	names := []string{"alice", "bob", "charlie"}
	fmt.Println("Contains bob:", contains(names, "bob"))
	fmt.Println("Contains dave:", contains(names, "dave"))

	ids := []int{10, 20, 30}
	fmt.Println("Contains 20:", contains(ids, 20))

	// Test 3: Minimum - now type-safe
	fmt.Println("\n--- Minimum ---")
	fmt.Println("Min(3, 7):", minimum(3, 7))
	fmt.Println("Min(3.14, 2.71):", minimum(3.14, 2.71))
	fmt.Println("Min(apple, banana):", minimum("apple", "banana"))

	// Test 4: Map - now type-safe, no type assertions needed
	fmt.Println("\n--- Map ---")
	words := []string{"hello", "world", "go"}
	upper := mapSlice(words, func(v string) string {
		return strings.ToUpper(v)
	})
	fmt.Println("Uppercased:", upper)

	nums := []int{1, 2, 3, 4}
	doubled := mapSlice(nums, func(v int) int {
		return v * 2
	})
	fmt.Println("Doubled:", doubled)

	fmt.Println("\nGenerics refactoring complete!")
}
