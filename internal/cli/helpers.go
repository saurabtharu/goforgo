package cli

import (
	"fmt"

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
