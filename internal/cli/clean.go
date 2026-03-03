package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove build artifacts from exercise directories",
	Long: `Remove binaries, go.mod, and go.sum files created by running exercises.

These files are generated automatically when exercises are compiled and
can be safely removed at any time.

Examples:
  goforgo clean`,
	RunE: cleanArtifacts,
}

func cleanArtifacts(cmd *cobra.Command, args []string) error {
	cwd, err := GetWorkingDirectory()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	exercisesDir := filepath.Join(cwd, "exercises")
	if _, err := os.Stat(exercisesDir); os.IsNotExist(err) {
		return fmt.Errorf("exercises directory not found. Run 'goforgo init' first")
	}

	removed := 0

	err = filepath.Walk(exercisesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		shouldRemove := false

		switch {
		case name == "go.mod" || name == "go.sum":
			shouldRemove = true
		case !strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, ".toml"):
			// Binary — no extension, not a source/metadata file
			shouldRemove = true
		}

		if shouldRemove {
			if err := os.Remove(path); err != nil {
				fmt.Printf("  warning: could not remove %s: %v\n", path, err)
			} else {
				removed++
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk exercises directory: %w", err)
	}

	if removed == 0 {
		fmt.Println("Nothing to clean")
	} else {
		fmt.Printf("Removed %d build artifacts\n", removed)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
