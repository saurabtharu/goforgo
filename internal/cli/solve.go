//go:build dev

package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve <N or X-Y>",
	Short: "Copy solutions over exercises for quick testing",
	Long: `Copy solution files over exercise files for a range of exercises,
marking them as completed. Exercise numbers are 1-indexed based on the
sorted exercise list.

Examples:
  goforgo solve 5        # Solve exercise 5
  goforgo solve 3-7      # Solve exercises 3 through 7`,
	Args: cobra.ExactArgs(1),
	RunE: solveExercises,
}

func solveExercises(cmd *cobra.Command, args []string) error {
	em, _, err := loadExerciseManager()
	if err != nil {
		return err
	}

	exercises := em.GetExercises()
	total := len(exercises)

	// Parse range argument
	start, end, err := parseRange(args[0], total)
	if err != nil {
		return err
	}

	solved := 0
	for i := start; i <= end; i++ {
		ex := exercises[i]

		// Read solution file
		solution, err := os.ReadFile(ex.SolutionPath)
		if err != nil {
			fmt.Printf("Skip: %s (no solution: %v)\n", ex.Info.Name, err)
			continue
		}

		// Overwrite exercise with solution
		if err := os.WriteFile(ex.FilePath, solution, 0644); err != nil {
			fmt.Printf("Skip: %s (write failed: %v)\n", ex.Info.Name, err)
			continue
		}

		// Mark completed
		if err := em.MarkExerciseCompleted(ex.Info.Name); err != nil {
			fmt.Printf("Skip: %s (mark failed: %v)\n", ex.Info.Name, err)
			continue
		}

		fmt.Printf("Solved: %s\n", ex.Info.Name)
		solved++
	}

	fmt.Printf("\nSolved %d exercises\n", solved)
	return nil
}

// parseRange parses "N" or "X-Y" (1-indexed, inclusive) into 0-indexed start/end.
func parseRange(arg string, total int) (int, int, error) {
	if strings.Contains(arg, "-") {
		parts := strings.SplitN(arg, "-", 2)
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid range start: %s", parts[0])
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid range end: %s", parts[1])
		}
		if start < 1 || end < start || end > total {
			return 0, 0, fmt.Errorf("range %d-%d out of bounds (1-%d)", start, end, total)
		}
		return start - 1, end - 1, nil
	}

	n, err := strconv.Atoi(arg)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid exercise number: %s", arg)
	}
	if n < 1 || n > total {
		return 0, 0, fmt.Errorf("exercise %d out of bounds (1-%d)", n, total)
	}
	return n - 1, n - 1, nil
}

func init() {
	rootCmd.AddCommand(solveCmd)
}
