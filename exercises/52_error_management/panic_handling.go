package main

import (
	"fmt"
	"strconv"
)

// 100 Go Mistakes #48: Panicking
//
// This code panics on normal error conditions (empty input, invalid format,
// out-of-range values) and uses recover() to catch them. This is an
// anti-pattern — panic is for programmer errors, not expected failures.
//
// FIX:
// 1. Change parseConfig() to return (map[string]string, error)
// 2. Change parseAge() to return (int, error)
// 3. Handle errors in main() with if err != nil checks
// 4. Keep the panic in mustCompileTemplate() — nil template IS a programmer error
// 5. Use recover() ONLY around mustCompileTemplate to catch the programming error

func parseConfig(data string) map[string]string {
	if data == "" {
		// BUG: Panicking on empty input — this is a normal error condition.
		// Users might pass empty strings. Return an error instead.
		panic("empty config data")
	}
	config := make(map[string]string)
	config["app"] = data
	return config
}

func parseAge(input string) int {
	age, err := strconv.Atoi(input)
	if err != nil {
		// BUG: Panicking on invalid user input.
		// Invalid input is expected — return an error.
		panic(fmt.Sprintf("invalid age: %v", err))
	}
	if age < 0 || age > 150 {
		// BUG: Panicking on out-of-range value.
		// This is validation, not a programming error.
		panic(fmt.Sprintf("age %d out of range", age))
	}
	return age
}

// mustCompileTemplate panics when template is nil.
// This IS appropriate — a nil template is a programmer error (a bug),
// not an expected runtime condition. The "must" prefix signals this.
func mustCompileTemplate(name string, template *string) string {
	if template == nil {
		panic(fmt.Sprintf("template %q must not be nil", name))
	}
	return fmt.Sprintf("[compiled: %s]", *template)
}

func main() {
	fmt.Println("=== Panic Handling ===")

	// BUG: Using recover() as error handling — this is an anti-pattern.
	// recover() should be used to catch programming errors or protect
	// goroutines, not to handle expected failures.

	fmt.Println("--- Config Parsing ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		config := parseConfig("")
		fmt.Printf("Config: %v\n", config)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		config := parseConfig("myapp")
		fmt.Printf("Config: %v\n", config)
	}()

	fmt.Println("--- Age Parsing ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		age := parseAge("abc")
		fmt.Printf("Age: %d\n", age)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		age := parseAge("-5")
		fmt.Printf("Age: %d\n", age)
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		age := parseAge("25")
		fmt.Printf("Age: %d\n", age)
	}()

	fmt.Println("--- Template Compilation ---")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("RECOVERED: %v\n", r)
			}
		}()
		result := mustCompileTemplate("header", nil)
		fmt.Println(result)
	}()

	tmpl := "<h1>Hello</h1>"
	result := mustCompileTemplate("header", &tmpl)
	fmt.Printf("Template: %s\n", result)
}
