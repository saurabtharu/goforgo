package main

import (
	"fmt"
	"strconv"
)

// 100 Go Mistakes #48: Panicking (Solution)
//
// Fixed: Functions return errors for expected failure conditions.
// panic is reserved for programming errors (impossible states).
// recover() is used only to catch programming errors gracefully.

func parseConfig(data string) (map[string]string, error) {
	if data == "" {
		// Fixed: Return an error instead of panicking.
		// Empty input is a normal condition, not a programming bug.
		return nil, fmt.Errorf("empty config data")
	}
	config := make(map[string]string)
	config["app"] = data
	return config, nil
}

func parseAge(input string) (int, error) {
	age, err := strconv.Atoi(input)
	if err != nil {
		// Fixed: Return error with context instead of panicking.
		return 0, fmt.Errorf("invalid age %q: %w", input, err)
	}
	if age < 0 || age > 150 {
		// Fixed: Return error for validation failure.
		return 0, fmt.Errorf("age %d out of valid range [0, 150]", age)
	}
	return age, nil
}

// mustCompileTemplate still panics — nil template is a programmer error.
// The "must" prefix convention tells callers this will panic on misuse.
func mustCompileTemplate(name string, template *string) string {
	if template == nil {
		panic(fmt.Sprintf("template %q must not be nil", name))
	}
	return fmt.Sprintf("[compiled: %s]", *template)
}

func main() {
	fmt.Println("=== Panic Handling ===")

	// Fixed: Handle errors with normal if err != nil checks.
	fmt.Println("--- Config Parsing ---")
	_, err := parseConfig("")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	config, err := parseConfig("myapp")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Config: %v\n", config)
	}

	fmt.Println("--- Age Parsing ---")
	_, err = parseAge("abc")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	_, err = parseAge("-5")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}

	age, err := parseAge("25")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Valid age: %d\n", age)
	}

	fmt.Println("--- Template Compilation ---")
	// recover() is appropriate here — catching a programming error
	// to log it gracefully (e.g., in a web server request handler).
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Caught programming error: %v\n", r)
			}
		}()
		mustCompileTemplate("header", nil)
	}()

	tmpl := "<h1>Hello</h1>"
	result := mustCompileTemplate("header", &tmpl)
	fmt.Printf("Template: %s\n", result)
}
