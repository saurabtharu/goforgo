package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	selfUpdateYes   bool
	selfUpdateCheck bool
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Check for and install the latest GoForGo version",
	Long: `Checks the latest GitHub tag and optionally updates the installed goforgo binary.

By default, this command asks for confirmation before updating.
Use --yes to skip the prompt.
Use --check to only check and print status.`,
	RunE: runSelfUpdate,
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	latest, isNewer, err := checkForUpdate(version, defaultTagsURLForTests, nil)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !isNewer {
		fmt.Fprintf(cmd.OutOrStdout(), "✅ goforgo is up to date (%s)\n", version)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "🔔 Update available: %s (current: %s)\n", latest, version)
	fmt.Fprintf(cmd.OutOrStdout(), "   Will run: %s\n", updateInstallCmd)

	if selfUpdateCheck {
		return nil
	}

	if !selfUpdateYes {
		ok, err := askForConfirmation(cmd)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "Update cancelled.")
			return nil
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Updating goforgo...")
	goCmd := exec.Command("go", "install", "github.com/stonecharioteer/goforgo/cmd/goforgo@latest")
	goCmd.Stdout = cmd.OutOrStdout()
	goCmd.Stderr = cmd.ErrOrStderr()

	if err := goCmd.Run(); err != nil {
		return fmt.Errorf("self-update failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "✅ Update complete. Run 'goforgo --version' to verify.")
	return nil
}

func askForConfirmation(cmd *cobra.Command) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), "Proceed with update? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
	selfUpdateCmd.Flags().BoolVarP(&selfUpdateYes, "yes", "y", false, "run update without confirmation prompt")
	selfUpdateCmd.Flags().BoolVar(&selfUpdateCheck, "check", false, "only check for updates, do not install")
}
