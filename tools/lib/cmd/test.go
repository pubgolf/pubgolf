package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

func init() {
	testCmd.AddCommand(testE2ECmd)
	testE2ECmd.PersistentFlags().Bool("watch", false, "Watch the input directory and automatically restart the e2e test")

	rootCmd.AddCommand(testCmd)

	testCmd.PersistentFlags().BoolP("coverage", "c", false, "Generate and display a coverage profile")
	testCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose test output")
	testCmd.PersistentFlags().Bool("local", false, "Run tests in a mode that disables external network dependencies")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all automated unit/integration tests",
	Run: func(cmd *cobra.Command, _ []string) {
		coverageDir := filepath.FromSlash("./data/go-test-coverage")
		coverageFile := filepath.Join(coverageDir, "api.cover")

		coverageFlag, err := cmd.Flags().GetBool("coverage")
		guard(err, "check '--coverage' flag")

		verboseFlag, err := cmd.Flags().GetBool("verbose")
		guard(err, "check '--verbose' flag")

		localFlag, err := cmd.Flags().GetBool("local")
		guard(err, "check '--local' flag")

		args := []string{
			"run",
			"--project", config.ServerBinName,
			"--config", "test",
		}

		if localFlag {
			args = append(args, []string{
				"--no-check-version",
				"--fallback-only",
			}...)
		}

		args = append(args, []string{
			"--",
			"go",
			"test",
			filepath.FromSlash("./api/..."),
		}...)

		if localFlag {
			args = append(args, "-shared-postgres=true")
		}

		if verboseFlag {
			args = append(args, "-v")
		}

		if coverageFlag {
			args = append(args,
				"-coverprofile", coverageFile,
			)

			guard(os.RemoveAll(coverageDir), "clean old output dir: %w")
			guard(os.MkdirAll(coverageDir, 0o755), "make new output dir: %w")
		}

		tester := exec.Command("doppler", args...)

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

var testE2ECmd = &cobra.Command{
	Use:   "e2e",
	Short: "Run all automated e2e tests",
	Run: func(cmd *cobra.Command, _ []string) {
		watchFlag, err := cmd.Flags().GetBool("watch")
		guard(err, "check '--watch' flag")

		localFlag, err := cmd.Flags().GetBool("local")
		guard(err, "check '--local' flag")

		// Start initial process.
		stopFn := runE2ETest(!watchFlag, localFlag)

		// Launch watcher, if applicable.
		if watchFlag {
			go func() {
				watch("api", "e2e test", func(_ watcher.Event) {
					stopFn()
					stopFn = runE2ETest(false, localFlag)
				})
			}()
		}

		// Hold process open
		<-shuttingDown
	},
}

func runE2ETest(stopOnExit, localOnly bool) func() {
	args := []string{
		"run",
		"--project", config.ServerBinName,
		"--config", "test",
	}

	if localOnly {
		args = append(args, []string{
			"--no-check-version",
			"--fallback-only",
		}...)
	}

	args = append(args, []string{
		"--",
		"go",
		"test",
		filepath.FromSlash("./api/internal/e2e"),
		"-v",
		"-e2e=true",
	}...)

	if localOnly {
		args = append(args, []string{
			"-shared-postgres=true",
		}...)
	}

	tester := exec.Command("doppler", args...)

	tester.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	tester.Stdout = os.Stdout
	tester.Stderr = os.Stderr
	tester.Stdin = os.Stdin

	log.Println("Starting e2e test run...")

	err := tester.Start()
	if err != nil {
		// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 { //nolint:errorlint // Casting to extract data.
			guard(err, "execute `go test ...` command")
		}
	}

	if stopOnExit {
		go func() {
			defer close(beginShutdown)

			err := tester.Wait()
			if err != nil {
				// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
				if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 { //nolint:errorlint // Casting to extract data.
					guard(err, "execute `go test ...` command")
				}
			}
		}()
	}

	return func() {
		pgid, err := syscall.Getpgid(tester.Process.Pid)
		if err != nil {
			// Previous run has already finished.
			return
		}

		guard(syscall.Kill(-pgid, syscall.SIGINT), "send SIGINT to running process")
	}
}
