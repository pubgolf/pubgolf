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
	Run: func(cmd *cobra.Command, args []string) {
		var coverageDir = filepath.FromSlash("./data/go-test-coverage")
		var coverageFile = filepath.Join(coverageDir, "api.cover")

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
			guard(os.MkdirAll(coverageDir, 0755), "make new output dir: %w")
		}

		tester := exec.Command("go", testArgs...)

		tester.Stdout = os.Stdout
		tester.Stderr = os.Stderr
		tester.Stdin = os.Stdin

		guard(tester.Run(), "execute `go test ...` command")

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
