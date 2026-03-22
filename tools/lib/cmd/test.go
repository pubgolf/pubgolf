package cmd

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

		err = runner.Run(cmd.Context(), Cmd{
			Name: "doppler",
			Args: args,
		})
		if err != nil {
			// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
			exitErr, ok := err.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
			if !ok || exitErr.ExitCode() != 1 {
				guard(err, "execute `go test ...` command")
			}
		}

		if coverageFlag {
			guard(runner.Run(cmd.Context(), Cmd{
				Name: "go",
				Args: []string{"tool", "cover", "-html", coverageFile},
			}), "execute `go tool cover ...` command")
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
		proc := runE2ETest(cmd.Context(), runner, !watchFlag, localFlag)

		// Launch watcher, if applicable.
		if watchFlag {
			go func() {
				watch("api", "e2e test", func(_ watcher.Event) {
					proc.Stop()
					proc = runE2ETest(context.Background(), runner, false, localFlag)
				})
			}()
		}

		// Hold process open
		<-shuttingDown
	},
}

func runE2ETest(ctx context.Context, r Runner, stopOnExit, localOnly bool) Process { //nolint:ireturn // Returns Process interface by design.
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
			"-doppler-local=true",
			"-disable-telemetry=true",
		}...)
	}

	log.Println("Starting e2e test run...")

	proc, err := r.Start(ctx, Cmd{
		Name: "doppler",
		Args: args,
	})
	if err != nil {
		// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
		exitErr, ok := err.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
		if !ok || exitErr.ExitCode() != 1 {
			guard(err, "execute `go test ...` command")
		}
	}

	if stopOnExit {
		go func() {
			defer triggerShutdown()

			err := proc.Wait()
			if err != nil {
				// Panic on error, unless the exit code is 1, in which case it just means our test suite failed.
				exitErr, ok := err.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
				if !ok || exitErr.ExitCode() != 1 {
					guard(err, "execute `go test ...` command")
				}
			}
		}()
	}

	return proc
}
