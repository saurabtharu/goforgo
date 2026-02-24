package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stonecharioteer/goforgo/internal/exercise"
)

// loadExerciseManager creates and initializes an ExerciseManager for the current
// working directory. Returns the manager and the working directory path.
func loadExerciseManager() (*exercise.ExerciseManager, string, error) {
	cwd, err := GetWorkingDirectory()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get working directory: %w", err)
	}

	em := exercise.NewExerciseManager(cwd)
	if err := em.LoadExercises(); err != nil {
		return nil, cwd, err
	}

	return em, cwd, nil
}

// findExerciseSourceDirs locates the exercise and solution source directories
// relative to the executable or the current working directory.
func findExerciseSourceDirs() (exercisesDir, solutionsDir string, err error) {
	execPath, execErr := os.Executable()
	if execErr != nil {
		execPath = ""
	}

	var possiblePaths []string
	if execPath != "" {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths,
			filepath.Join(execDir, "exercises"),
			filepath.Join(execDir, "..", "exercises"),
			filepath.Join(execDir, "..", "..", "exercises"),
		)
	}

	possiblePaths = append(possiblePaths,
		"exercises",
		"../exercises",
		"../../exercises",
	)

	for _, path := range possiblePaths {
		if _, statErr := os.Stat(path); statErr == nil {
			exercisesDir = path
			solutionsDir = strings.Replace(path, "exercises", "solutions", 1)
			return exercisesDir, solutionsDir, nil
		}
	}

	return "", "", fmt.Errorf("no source exercises found")
}
