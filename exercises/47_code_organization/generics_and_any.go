// generics_and_any.go
// Replace unsafe any/interface{} usage with proper type constraints
//
// Using any (or interface{}) throws away type safety. The compiler can't
// catch type mismatches, so bugs surface at runtime instead of compile time.
// Go 1.18+ generics let you express type constraints that the compiler enforces.
//
// This code uses any where type parameters would catch bugs at compile time.
// Refactor to use proper generic type constraints.

package main

import (
	"fmt"
	"strings"
)

// BUG 1: This function uses any but only works with numeric types.
// It will panic at runtime if you pass a string.
// TODO: Replace any with a numeric type constraint.
func sum(values []any) any {
	var total float64
	for _, v := range values {
		switch n := v.(type) {
		case int:
			total += float64(n)
		case float64:
			total += n
		default:
			panic(fmt.Sprintf("unsupported type: %T", v))
		}
	}
	return total
}

// BUG 2: This uses any for a "contains" check, but the comparison
// with == only works for comparable types. A slice passed here panics.
// TODO: Use a generic with the comparable constraint.
func contains(slice []any, target any) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// BUG 3: This min function works on any but needs runtime type assertions.
// TODO: Rewrite using a type constraint that supports the < operator.
// Hint: Use constraints like "int | float64 | string" or cmp.Ordered
// if you don't want to define your own.
type Ordered interface {
	~int | ~float64 | ~string
}

func minimum(a, b any) any {
	switch va := a.(type) {
	case int:
		vb := b.(int)
		if va < vb {
			return va
		}
		return vb
	case float64:
		vb := b.(float64)
		if va < vb {
			return va
		}
		return vb
	case string:
		vb := b.(string)
		if va < vb {
			return va
		}
		return vb
	default:
		panic(fmt.Sprintf("unsupported type: %T", a))
	}
}

// BUG 4: This map function uses any for both input and output.
// The caller has to do type assertions on every result.
// TODO: Make this a generic function with type parameters.
func mapSlice(input []any, fn func(any) any) []any {
	result := make([]any, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func main() {
	fmt.Println("=== Generics and Any ===")

	// Test 1: Sum
	fmt.Println("\n--- Sum ---")
	ints := []any{1, 2, 3, 4, 5}
	fmt.Println("Sum of ints:", sum(ints))

	floats := []any{1.5, 2.5, 3.5}
	fmt.Println("Sum of floats:", sum(floats))

	// Test 2: Contains
	fmt.Println("\n--- Contains ---")
	names := []any{"alice", "bob", "charlie"}
	fmt.Println("Contains bob:", contains(names, "bob"))
	fmt.Println("Contains dave:", contains(names, "dave"))

	ids := []any{10, 20, 30}
	fmt.Println("Contains 20:", contains(ids, 20))

	// Test 3: Minimum
	fmt.Println("\n--- Minimum ---")
	fmt.Println("Min(3, 7):", minimum(3, 7))
	fmt.Println("Min(3.14, 2.71):", minimum(3.14, 2.71))
	fmt.Println("Min(apple, banana):", minimum("apple", "banana"))

	// Test 4: Map
	fmt.Println("\n--- Map ---")
	words := []any{"hello", "world", "go"}
	upper := mapSlice(words, func(v any) any {
		return strings.ToUpper(v.(string))
	})
	fmt.Println("Uppercased:", upper)

	nums := []any{1, 2, 3, 4}
	doubled := mapSlice(nums, func(v any) any {
		return v.(int) * 2
	})
	fmt.Println("Doubled:", doubled)

	fmt.Println("\nGenerics refactoring complete!")
}
