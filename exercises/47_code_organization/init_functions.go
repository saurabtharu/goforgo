// init_functions.go
// Understand init() function pitfalls and refactor to explicit initialization
//
// Go's init() functions run automatically before main(), in declaration order
// within a file and by import order across packages. This makes them:
// - Hard to test (no way to skip or mock)
// - Surprising (side effects happen before main)
// - Order-dependent (fragile when refactoring)
//
// This program uses init() for global state setup. Refactor it to use
// explicit initialization functions called from main() instead.

package main

import (
	"fmt"
	"strings"
)

// Global state set by init functions - this is the anti-pattern.
// TODO: Remove these init() functions and replace them with explicit
// initialization functions. Call them from main() in the right order.

var appName string
var appVersion string
var features []string
var config map[string]string

// BUG: This init() sets global state silently.
// TODO: Convert to an explicit function like initApp() (string, string)
func init() {
	appName = "GoForGo"
	appVersion = "1.0.0"
	fmt.Println("init 1: App identity set")
}

// BUG: This init() depends on the first init() having run.
// TODO: Convert to an explicit function like initFeatures() []string
func init() {
	features = []string{
		strings.ToLower(appName) + "-exercises",
		strings.ToLower(appName) + "-tui",
		strings.ToLower(appName) + "-hints",
	}
	fmt.Println("init 2: Features configured")
}

// BUG: This init() depends on both previous init() functions.
// TODO: Convert to an explicit function like initConfig() map[string]string
func init() {
	config = map[string]string{
		"app":      appName,
		"version":  appVersion,
		"features": strings.Join(features, ","),
		"mode":     "development",
	}
	fmt.Println("init 3: Config built")
}

func printStatus() {
	fmt.Printf("App: %s v%s\n", appName, appVersion)
	fmt.Printf("Features: %s\n", strings.Join(features, ", "))
	fmt.Println("Config:")
	// Print in deterministic order
	for _, key := range []string{"app", "version", "features", "mode"} {
		fmt.Printf("  %s = %s\n", key, config[key])
	}
}

func main() {
	fmt.Println("\n=== Init Functions ===")
	fmt.Println("(Notice the init messages printed BEFORE main starts)")
	fmt.Println()

	printStatus()

	// TODO: After refactoring, the output should be:
	// === Init Functions ===
	// Initializing app identity...
	// Initializing features...
	// Initializing config...
	//
	// App: GoForGo v1.0.0
	// Features: goforgo-exercises, goforgo-tui, goforgo-hints
	// Config:
	//   app = GoForGo
	//   version = 1.0.0
	//   features = goforgo-exercises,goforgo-tui,goforgo-hints
	//   mode = development
	//
	// Init refactoring complete!

	fmt.Println("\nInit refactoring complete!")
}
