package main

import (
	"fmt"
	"reflect"
)

// 100 Go Mistakes #29: Comparing Values Incorrectly (Solution)
//
// Fixed:
// 1. Use reflect.DeepEqual for structs with slices/maps
// 2. Custom comparison for nil-vs-empty equivalence
// 3. Manual loop comparison for simple slices (faster than reflect)

type User struct {
	Name  string
	Email string
	Age   int
}

type Team struct {
	Name    string
	Members []string
	Scores  map[string]int
}

type Config struct {
	Version int
	Tags    []string
	Options map[string]string
}

func compareUsers(a, b User) bool {
	return a == b
}

func compareTeams(a, b Team) bool {
	// Fixed: Use reflect.DeepEqual for structs with uncomparable fields.
	// This handles slice and map comparison correctly.
	return reflect.DeepEqual(a, b)
}

// Fixed: Custom comparison that treats nil and empty as equivalent.
func compareConfigs(a, b Config) bool {
	if a.Version != b.Version {
		return false
	}

	// Compare Tags: treat nil and empty as equivalent
	if len(a.Tags) != len(b.Tags) {
		return false
	}
	for i := range a.Tags {
		if a.Tags[i] != b.Tags[i] {
			return false
		}
	}

	// Compare Options: treat nil and empty as equivalent
	if len(a.Options) != len(b.Options) {
		return false
	}
	for k, v := range a.Options {
		if bv, ok := b.Options[k]; !ok || v != bv {
			return false
		}
	}

	return true
}

// Fixed: Manual loop comparison for slices.
// This is faster than reflect.DeepEqual and more explicit.
func compareSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("=== Comparing Values ===")
	fmt.Println()

	fmt.Println("--- Comparable Structs ---")
	u1 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	u2 := User{Name: "Alice", Email: "alice@example.com", Age: 30}
	u3 := User{Name: "Bob", Email: "bob@example.com", Age: 25}
	fmt.Printf("u1 == u2: %v (expected: true)\n", compareUsers(u1, u2))
	fmt.Printf("u1 == u3: %v (expected: false)\n", compareUsers(u1, u3))
	fmt.Println()

	fmt.Println("--- Structs with Slices ---")
	t1 := Team{Name: "Alpha", Members: []string{"Alice", "Bob"}, Scores: map[string]int{"Alice": 10}}
	t2 := Team{Name: "Alpha", Members: []string{"Alice", "Bob"}, Scores: map[string]int{"Alice": 10}}
	t3 := Team{Name: "Alpha", Members: []string{"Alice", "Charlie"}, Scores: map[string]int{"Alice": 10}}
	fmt.Printf("t1 == t2: %v (expected: true)\n", compareTeams(t1, t2))
	fmt.Printf("t1 == t3: %v (expected: false)\n", compareTeams(t1, t3))
	fmt.Println()

	fmt.Println("--- Nil vs Empty Equivalence ---")
	c1 := Config{Version: 1, Tags: nil, Options: nil}
	c2 := Config{Version: 1, Tags: []string{}, Options: map[string]string{}}
	c3 := Config{Version: 2, Tags: nil, Options: nil}
	fmt.Printf("c1 == c2 (nil vs empty): %v (expected: true)\n", compareConfigs(c1, c2))
	fmt.Printf("c1 == c3 (different version): %v (expected: false)\n", compareConfigs(c1, c3))
	fmt.Println()

	fmt.Println("--- Slice Comparison ---")
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{1, 2, 4}
	var s4 []int
	fmt.Printf("[1,2,3] == [1,2,3]: %v (expected: true)\n", compareSlices(s1, s2))
	fmt.Printf("[1,2,3] == [1,2,4]: %v (expected: false)\n", compareSlices(s1, s3))
	fmt.Printf("[1,2,3] == nil:     %v (expected: false)\n", compareSlices(s1, s4))
}
