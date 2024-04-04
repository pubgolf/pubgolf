package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.PersistentFlags().BoolP("coverage", "c", false, "Generate and display a coverage profile")
	testCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose test output")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all automated tests",
	Run: func(cmd *cobra.Command, _ []string) {
		coverageDir := filepath.FromSlash("./data/go-test-coverage")
		coverageFile := filepath.Join(coverageDir, "api.cover")

		coverageFlag, err := cmd.Flags().GetBool("coverage")
		guard(err, "check '--coverage' flag")

		verboseFlag, err := cmd.Flags().GetBool("verbose")
		guard(err, "check '--verbose' flag")

		testArgs := []string{
			"test",

			filepath.FromSlash("./api/..."),
		}

		if verboseFlag {
			testArgs = append(testArgs, "-v")
		}

		if coverageFlag {
			testArgs = append(testArgs,
				"-coverprofile", coverageFile,
			)

			guard(os.RemoveAll(coverageDir), "clean old output dir: %w")
			guard(os.MkdirAll(coverageDir, 0o755), "make new output dir: %w")
		}

		tester := exec.Command("go", testArgs...)

		tester.Stdout = os.Stdout
		tester.Stderr = os.Stderr
		tester.Stdin = os.Stdin

		err = tester.Run()
		if err != nil {
			// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
			if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 { //nolint:errorlint // Casting to extract data.
				guard(err, "execute `go test ...` command")
			}
		}

		if coverageFlag {
			cover := exec.Command("go",
				"tool", "cover",
				"-html", coverageFile,
			)

			cover.Stdout = os.Stdout
			cover.Stderr = os.Stderr
			cover.Stdin = os.Stdin

			guard(cover.Run(), "execute `go tool cover ...` command")
		}
	},
}
