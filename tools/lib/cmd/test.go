package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		coverageDir := filepath.Join(projectRoot, worktreeDataDir(cmd.Context(), "data/go-test-coverage"))
		coverageFile := filepath.Join(coverageDir, "api.cover")

		coverageFlag, err := cmd.Flags().GetBool("coverage")
		classifyAndExit(fmtErr(err, "check '--coverage' flag"))

		verboseFlag, err := cmd.Flags().GetBool("verbose")
		classifyAndExit(fmtErr(err, "check '--verbose' flag"))

		localFlag, err := cmd.Flags().GetBool("local")
		classifyAndExit(fmtErr(err, "check '--local' flag"))

		env, envErr := envProvider.Env(cmd.Context(), config.ServerBinName, "test")
		classifyAndExit(fmtErr(envErr, "fetch test environment"))

		// Inject worktree isolation env vars.
		slug, slugErr := worktreeSlug(cmd.Context())
		classifyAndExit(fmtErr(slugErr, "determine worktree slug"))

		if slug != "" {
			env = append(env, "PUBGOLF_WORKTREE_SLUG="+slug)
		}

		offset, offsetErr := worktreePortOffset(cmd.Context())
		classifyAndExit(fmtErr(offsetErr, "compute port offset"))

		if localFlag && offset > 0 {
			env = replaceSharedDBPort(env, offset)
		}

		// Pre-flight infrastructure check for shared-postgres mode.
		if localFlag {
			classifyAndExit(preflight(cmd.Context(), offset))
		}

		args := []string{"test", filepath.FromSlash("./api/...")}

		if localFlag {
			args = append(args, "-shared-postgres=true")
		}

		if verboseFlag {
			args = append(args, "-v")
		}

		if coverageFlag {
			args = append(args, "-coverprofile", coverageFile)

			classifyAndExit(fmtErr(os.RemoveAll(coverageDir), "clean old output dir"))
			classifyAndExit(fmtErr(os.MkdirAll(coverageDir, 0o755), "make new output dir"))
		}

		err = runner.Run(cmd.Context(), Cmd{
			Name: "go",
			Args: args,
			Env:  env,
		})
		if err != nil {
			// Allow exit code 1 through — it means the test suite failed (not infrastructure).
			exitErr, ok := err.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
			if !ok || exitErr.ExitCode() != 1 {
				classifyAndExit(fmtErr(err, "execute `go test ...` command"))
			}
		}

		if coverageFlag {
			classifyAndExit(fmtErr(runner.Run(cmd.Context(), Cmd{
				Name: "go",
				Args: []string{"tool", "cover", "-html", coverageFile},
			}), "execute `go tool cover ...` command"))
		}
	},
}

var testE2ECmd = &cobra.Command{
	Use:   "e2e",
	Short: "Run all automated e2e tests",
	Run: func(cmd *cobra.Command, _ []string) {
		watchFlag, err := cmd.Flags().GetBool("watch")
		classifyAndExit(fmtErr(err, "check '--watch' flag"))

		localFlag, err := cmd.Flags().GetBool("local")
		classifyAndExit(fmtErr(err, "check '--local' flag"))

		// Start initial process.
		proc := runE2ETest(cmd.Context(), runner, envProvider, !watchFlag, localFlag)

		// Launch watcher, if applicable.
		if watchFlag {
			go func() {
				watch("api", "e2e test", func(_ watcher.Event) {
					proc.Stop()
					proc = runE2ETest(context.Background(), runner, envProvider, false, localFlag)
				})
			}()
		}

		// Hold process open
		<-shuttingDown
	},
}

func runE2ETest(ctx context.Context, r Runner, ep EnvProvider, stopOnExit, localOnly bool) Process { //nolint:ireturn // Returns Process interface by design.
	env, err := ep.Env(ctx, config.ServerBinName, "test") //nolint:ireturn // False positive: linter reports ireturn on body line.
	classifyAndExit(fmtErr(err, "fetch e2e test environment"))

	// Inject worktree isolation env vars.
	slug, slugErr := worktreeSlug(ctx)
	classifyAndExit(fmtErr(slugErr, "determine worktree slug"))

	if slug != "" {
		env = append(env, "PUBGOLF_WORKTREE_SLUG="+slug)
	}

	offset, offsetErr := worktreePortOffset(ctx)
	classifyAndExit(fmtErr(offsetErr, "compute port offset"))

	if localOnly && offset > 0 {
		env = replaceSharedDBPort(env, offset)
	}

	args := []string{
		"test",
		filepath.FromSlash("./api/internal/e2e"),
		"-v",
		"-e2e=true",
	}

	if localOnly {
		args = append(args, []string{
			"-shared-postgres=true",
			"-doppler-local=true",
			"-disable-telemetry=true",
		}...)
	}

	log.Println("Starting e2e test run...")

	proc, startErr := r.Start(ctx, Cmd{
		Name: "go",
		Args: args,
		Env:  env,
	})
	if startErr != nil {
		// Allow exit code 1 through — it means the test suite failed (not infrastructure).
		exitErr, ok := startErr.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
		if !ok || exitErr.ExitCode() != 1 {
			classifyAndExit(fmtErr(startErr, "execute `go test ...` command"))
		}
	}

	if stopOnExit {
		go func() {
			defer triggerShutdown()

			waitErr := proc.Wait()
			if waitErr != nil {
				// Allow exit code 1 through — it means the test suite failed (not infrastructure).
				exitErr, ok := waitErr.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
				if !ok || exitErr.ExitCode() != 1 {
					classifyAndExit(fmtErr(waitErr, "execute `go test ...` command"))
				}
			}
		}()
	}

	return proc
}

// errDBNotRunning is returned when the pre-flight health check cannot reach the database.
var errDBNotRunning = errors.New("DB container is not reachable")

// preflight checks that the database container is reachable before running tests.
func preflight(ctx context.Context, offset int) error {
	dbPort := 5432 + offset

	dialer := &net.Dialer{Timeout: 2 * time.Second}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%d", dbPort))
	if err != nil {
		return fmt.Errorf("%w on port %d.\n"+
			"  Run 'pubgolf-devctrl run bg' to start it.\n"+
			"  This is an infrastructure issue, not a code problem", errDBNotRunning, dbPort)
	}

	conn.Close()

	return nil
}

// replaceSharedDBPort finds PUBGOLF_SHARED_DB_URL in the env slice, parses it,
// replaces the port with the offset port, and returns the updated env slice.
// If the variable is not present or cannot be parsed, the slice is returned unchanged.
func replaceSharedDBPort(env []string, offset int) []string {
	const key = "PUBGOLF_SHARED_DB_URL="

	for i, v := range env {
		if !strings.HasPrefix(v, key) {
			continue
		}

		raw := strings.TrimPrefix(v, key)

		u, err := url.Parse(raw)
		if err != nil {
			return env
		}

		u.Host = net.JoinHostPort(u.Hostname(), strconv.Itoa(5432+offset))
		env[i] = key + u.String()

		return env
	}

	return env
}
