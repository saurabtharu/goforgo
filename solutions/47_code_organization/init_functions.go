// init_functions.go - SOLUTION
// Refactored from init() functions to explicit initialization.

package main

import (
	"fmt"
	"strings"
)

var appName string
var appVersion string
var features []string
var config map[string]string

// Fixed: Explicit initialization function instead of init().
// Returns the values so the caller controls when and how they're set.
func initApp() (string, string) {
	fmt.Println("Initializing app identity...")
	return "GoForGo", "1.0.0"
}

// Fixed: Takes dependencies as parameters instead of relying on global state.
func initFeatures(name string) []string {
	fmt.Println("Initializing features...")
	return []string{
		strings.ToLower(name) + "-exercises",
		strings.ToLower(name) + "-tui",
		strings.ToLower(name) + "-hints",
	}
}

// Fixed: All dependencies are explicit parameters.
func initConfig(name, version string, feats []string) map[string]string {
	fmt.Println("Initializing config...")
	return map[string]string{
		"app":      name,
		"version":  version,
		"features": strings.Join(feats, ","),
		"mode":     "development",
	}
}

func printStatus() {
	fmt.Printf("App: %s v%s\n", appName, appVersion)
	fmt.Printf("Features: %s\n", strings.Join(features, ", "))
	fmt.Println("Config:")
	for _, key := range []string{"app", "version", "features", "mode"} {
		fmt.Printf("  %s = %s\n", key, config[key])
	}
}

func main() {
	fmt.Println("=== Init Functions ===")

	// Fixed: Explicit initialization in the right order, visible in main().
	appName, appVersion = initApp()
	features = initFeatures(appName)
	config = initConfig(appName, appVersion, features)

	fmt.Println()
	printStatus()

	fmt.Println("\nInit refactoring complete!")
}
