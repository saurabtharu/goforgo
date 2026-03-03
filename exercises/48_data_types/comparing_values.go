package main

import (
	"fmt"
	"reflect"
)

// 100 Go Mistakes #29: Comparing Values Incorrectly
//
// This exercise covers comparison pitfalls in Go:
// 1. == cannot compare slices, maps, or structs with uncomparable fields
// 2. reflect.DeepEqual has gotchas (nil vs empty slice, performance)
// 3. When to use custom comparison functions
//
// FIX the comparison functions so all tests pass.

// User is a struct with comparable fields only.
type User struct {
	Name  string
	Email string
	Age   int
}

// Team has a slice field, making it uncomparable with ==.
type Team struct {
	Name    string
	Members []string
	Scores  map[string]int
}

// Config has a nested uncomparable field.
type Config struct {
	Version int
	Tags    []string
	Options map[string]string
}

// compareUsers tries to compare two Users.
func compareUsers(a, b User) bool {
	// This works because all fields of User are comparable.
	return a == b
}

// compareTeams tries to compare two Teams.
func compareTeams(a, b Team) bool {
	// FIX: Slices and maps are not comparable with ==.
	// This won't even compile!
	// Use reflect.DeepEqual or write a custom comparison.
	// But beware: reflect.DeepEqual treats nil and empty slices as different!
	return reflect.DeepEqual(a, b)
}

// compareConfigs compares two Configs, treating nil and empty
// slices/maps as equivalent (unlike reflect.DeepEqual).
func compareConfigs(a, b Config) bool {
	// FIX: reflect.DeepEqual considers nil slice != empty slice,
	// and nil map != empty map. For Config comparison, we want
	// nil and empty to be treated the same.
	// Write a custom comparison function instead.
	return reflect.DeepEqual(a, b)
}

// compareSlices demonstrates slice comparison.
func compareSlices(a, b []int) bool {
	// FIX: Slices cannot be compared with ==.
	// reflect.DeepEqual works but is slow.
	// For simple slices, a manual loop is faster and clearer.
	return reflect.DeepEqual(a, b)
}

func main() {
	fmt.Println("=== Comparing Values ===")
	fmt.Println()

	// Test 1: Struct comparison (works with ==)
	fmt.Println("--- Comparable Structs ---")
	u1 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	u2 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	u3 := User{Name: "Bob", Email: "bob@example.com", Age: 25}
	fmt.Printf("u1 == u2: %v (expected: true)\n", compareUsers(u1, u2))
	fmt.Printf("u1 == u3: %v (expected: false)\n", compareUsers(u1, u3))
	fmt.Println()

	// Test 2: Struct with slices
	fmt.Println("--- Structs with Slices ---")
	t1 := Team{Name: "Alpha", Members: []string{"Alice", "Bob"}, Scores: map[string]int{"Alice": 10}}
	t2 := Team{Name: "Alpha", Members: []string{"Alice", "Bob"}, Scores: map[string]int{"Alice": 10}}
	t3 := Team{Name: "Alpha", Members: []string{"Alice", "Charlie"}, Scores: map[string]int{"Alice": 10}}
	fmt.Printf("t1 == t2: %v (expected: true)\n", compareTeams(t1, t2))
	fmt.Printf("t1 == t3: %v (expected: false)\n", compareTeams(t1, t3))
	fmt.Println()

	// Test 3: nil vs empty equivalence
	fmt.Println("--- Nil vs Empty Equivalence ---")
	c1 := Config{Version: 1, Tags: nil, Options: nil}
	c2 := Config{Version: 1, Tags: []string{}, Options: map[string]string{}}
	c3 := Config{Version: 2, Tags: nil, Options: nil}
	fmt.Printf("c1 == c2 (nil vs empty): %v (expected: true)\n", compareConfigs(c1, c2))
	fmt.Printf("c1 == c3 (different version): %v (expected: false)\n", compareConfigs(c1, c3))
	fmt.Println()

	// Test 4: Slice comparison
	fmt.Println("--- Slice Comparison ---")
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}
	var s4 []int
	fmt.Printf("[1,2,3] == [1,2,3]: %v (expected: true)\n", compareSlices(s1, s2))
	fmt.Printf("[1,2,3] == [1,2,4]: %v (expected: false)\n", compareSlices(s1, s3))
	fmt.Printf("[1,2,3] == nil:     %v (expected: false)\n", compareSlices(s1, s4))
}
