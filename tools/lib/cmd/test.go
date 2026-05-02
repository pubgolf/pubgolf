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

	"github.com/spf13/cobra"
)

func init() {
	testCmd.AddCommand(testGoCmd)
	testCmd.AddCommand(testWebCmd)
	testCmd.AddCommand(testE2ECmd)
	testE2ECmd.AddCommand(testE2EAPICmd)
	testE2ECmd.AddCommand(testE2EWebCmd)

	rootCmd.AddCommand(testCmd)

	testCmd.PersistentFlags().BoolP("coverage", "c", false, "Generate and display a coverage profile (Go tests only)")
	testCmd.PersistentFlags().BoolP("verbose", "v", false, "Display verbose test output (Go tests only)")
	testCmd.PersistentFlags().Bool("local", false, "Run tests in a mode that disables external network dependencies (Go tests only)")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all automated unit/integration tests (Go + web)",
	Run: func(cmd *cobra.Command, _ []string) {
		opts, err := readTestGoFlags(cmd)
		classifyAndExit(err)

		classifyAndExit(runPar(cmd.Context(), runner,
			func(ctx context.Context, r Runner) error {
				return testGo(ctx, r, envProvider, opts)
			},
			testWeb,
		))
	},
}

type testGoOpts struct {
	coverage bool
	verbose  bool
	local    bool
}

var testGoCmd = &cobra.Command{
	Use:   "go",
	Short: "Run Go API unit/integration tests",
	Run: func(cmd *cobra.Command, _ []string) {
		opts, err := readTestGoFlags(cmd)
		classifyAndExit(err)
		classifyAndExit(testGo(cmd.Context(), runner, envProvider, opts))
	},
}

func readTestGoFlags(cmd *cobra.Command) (testGoOpts, error) {
	coverage, err := cmd.Flags().GetBool("coverage")
	if err != nil {
		return testGoOpts{}, fmtErr(err, "check '--coverage' flag")
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return testGoOpts{}, fmtErr(err, "check '--verbose' flag")
	}

	local, err := cmd.Flags().GetBool("local")
	if err != nil {
		return testGoOpts{}, fmtErr(err, "check '--local' flag")
	}

	return testGoOpts{coverage: coverage, verbose: verbose, local: local}, nil
}

func testGo(ctx context.Context, r Runner, ep EnvProvider, opts testGoOpts) error {
	coverageDir := filepath.Join(projectRoot, worktreeDataDir(ctx, "data/go-test-coverage"))
	coverageFile := filepath.Join(coverageDir, "api.cover")

	env, err := ep.Env(ctx, config.ServerBinName, "test")
	if err != nil {
		return fmtErr(err, "fetch test environment")
	}

	slug, err := worktreeSlug(ctx)
	if err != nil {
		return fmtErr(err, "determine worktree slug")
	}

	if slug != "" {
		env = append(env, "PUBGOLF_WORKTREE_SLUG="+slug)
	}

	offset, err := worktreePortOffset(ctx)
	if err != nil {
		return fmtErr(err, "compute port offset")
	}

	if opts.local && offset > 0 {
		env = replaceSharedDBPort(env, offset)
	}

	if opts.local {
		err := preflight(ctx, offset)
		if err != nil {
			return err
		}
	}

	args := []string{"test", filepath.FromSlash("./api/...")}

	if opts.local {
		args = append(args, "-shared-postgres=true")
	}

	if opts.verbose {
		args = append(args, "-v")
	}

	if opts.coverage {
		args = append(args, "-coverprofile", coverageFile)

		err := os.RemoveAll(coverageDir)
		if err != nil {
			return fmtErr(err, "clean old output dir")
		}

		err = os.MkdirAll(coverageDir, 0o755)
		if err != nil {
			return fmtErr(err, "make new output dir")
		}
	}

	err = r.Run(ctx, Cmd{
		Name: "go",
		Args: args,
		Env:  env,
	})
	if err != nil {
		// Surface all errors including exit-code-1 test failures —
		// classifyAndExit routes infra vs test via stderr substring matching,
		// not exit code.
		return fmtErr(err, "execute `go test ...` command")
	}

	if opts.coverage {
		err := r.Run(ctx, Cmd{
			Name: "go",
			Args: []string{"tool", "cover", "-html", coverageFile},
		})
		if err != nil {
			return fmtErr(err, "execute `go tool cover ...` command")
		}
	}

	return nil
}

var testWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Run vitest unit tests on web code",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(testWeb(cmd.Context(), runner))
	},
}

func testWeb(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: filepath.FromSlash("./node_modules/.bin/vitest"),
		Args: []string{"--run"},
		Dir:  "web-app",
	})
	if err != nil {
		return fmtErr(err, "run vitest unit tests")
	}

	return nil
}

var testE2ECmd = &cobra.Command{
	Use:   "e2e",
	Short: "Run all automated e2e tests",
	Run: func(cmd *cobra.Command, _ []string) {
		localFlag, err := cmd.Flags().GetBool("local")
		classifyAndExit(fmtErr(err, "check '--local' flag"))

		classifyAndExit(runPar(cmd.Context(), runner,
			func(ctx context.Context, r Runner) error {
				return testE2EAPI(ctx, r, envProvider, localFlag)
			},
			testE2EWeb,
		))
	},
}

var testE2EAPICmd = &cobra.Command{
	Use:   "api",
	Short: "Run Go API e2e tests",
	Run: func(cmd *cobra.Command, _ []string) {
		localFlag, err := cmd.Flags().GetBool("local")
		classifyAndExit(fmtErr(err, "check '--local' flag"))

		classifyAndExit(testE2EAPI(cmd.Context(), runner, envProvider, localFlag))
	},
}

var testE2EWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Run Playwright e2e tests on web code",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(testE2EWeb(cmd.Context(), runner))
	},
}

func testE2EAPI(ctx context.Context, r Runner, ep EnvProvider, localOnly bool) error {
	env, err := ep.Env(ctx, config.ServerBinName, "test")
	if err != nil {
		return fmtErr(err, "fetch e2e test environment")
	}

	// Inject worktree isolation env vars.
	slug, slugErr := worktreeSlug(ctx)
	if slugErr != nil {
		return fmtErr(slugErr, "determine worktree slug")
	}

	if slug != "" {
		env = append(env, "PUBGOLF_WORKTREE_SLUG="+slug)
	}

	offset, offsetErr := worktreePortOffset(ctx)
	if offsetErr != nil {
		return fmtErr(offsetErr, "compute port offset")
	}

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

	err = r.Run(ctx, Cmd{
		Name: "go",
		Args: args,
		Env:  env,
	})
	if err != nil {
		exitErr, ok := err.(*exec.ExitError) //nolint:errorlint // Casting to extract data.
		if !ok || exitErr.ExitCode() != 1 {
			return fmtErr(err, "execute API e2e tests")
		}
	}

	return nil
}

func testE2EWeb(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: filepath.FromSlash("./node_modules/.bin/playwright"),
		Args: []string{"test"},
		Dir:  "web-app",
	})
	if err != nil {
		return fmtErr(err, "run Playwright e2e tests")
	}

	return nil
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
