// Package cmd contains handlers for each CLI command.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	"golang.org/x/mod/sumdb/dirhash"
)

// errProjectRootNotFound is returned when go.mod cannot be found in any parent directory.
var errProjectRootNotFound = errors.New("go.mod not found in any parent directory")

// Internal plumbing.
var (
	// shuttingDown is a broadcast channel that closes to tell all processes to begin cleanup.
	shuttingDown = make(chan struct{})
	// shutdownOnce ensures triggerShutdown is called at most once.
	shutdownOnce sync.Once
)

// Package parameters.
var (
	// installedToolsHash holds the hash of the source code for the currently installed version of the binary.
	installedToolsHash string
	// config holds the provided configuration options.
	config CLIConfig
	// projectRoot holds the resolved absolute path to the project root directory.
	projectRoot string
	// runner is the package-level Runner used by all command handlers.
	runner Runner
	// envProvider is the package-level EnvProvider used by commands that need secrets.
	envProvider EnvProvider
)

// Execute is the entrypoint for calling the CLI.
func Execute(toolsDirHash string, c CLIConfig) {
	installedToolsHash = toolsDirHash

	c.setDefaults()
	config = c

	log.SetPrefix(fmt.Sprintf("[%s] ", config.CLIName))

	root, err := resolveProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s must be run from within the project directory: %v\n", config.CLIName, err)
		os.Exit(1)
	}

	projectRoot = root
	log.Printf("Resolved project root: %s", projectRoot)

	err = rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolveProjectRoot walks up from CWD looking for go.mod to find the project root.
func resolveProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		_, statErr := os.Stat(filepath.Join(dir, "go.mod"))
		if statErr == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errProjectRootNotFound
		}

		dir = parent
	}
}

// triggerShutdown closes the shuttingDown channel exactly once.
func triggerShutdown() {
	shutdownOnce.Do(func() { close(shuttingDown) })
}

var rootCmd = &cobra.Command{
	Use:   config.CLIName,
	Short: "DevCtrl is a task runner for local dev",
	Long:  `An opinionated task runner for personal use by @thedeerchild`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		// Initialize the runner based on the --dry-run flag.
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if dryRun {
			runner = &DryRunner{}
		} else {
			runner = ExecRunner{}
		}

		// Initialize the env provider.
		envProviderFlag, _ := cmd.Flags().GetString("env-provider")
		envProvider = newEnvProvider(envProviderFlag, dryRun, runner)

		// Skip the update warning for the update command itself and in dry-run mode.
		if !dryRun && cmd.CommandPath() != " update" {
			checkVersion()
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigCh
			triggerShutdown()
		}()
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("dry-run", false, "Print commands that would be executed without running them")
	rootCmd.PersistentFlags().String("env-provider", "auto", "Environment provider: auto, doppler, or env")
}

func checkVersion() {
	curToolsHash, err := dirhash.HashDir(filepath.Join(projectRoot, "tools"), "", dirhash.DefaultHash)
	classifyAndExit(fmtErr(err, "hash tools dir"))

	if installedToolsHash != curToolsHash {
		log.Printf(`The installed version of %[1]s is out of date. Run the following to update:

%[1]s update

`, config.CLIName)
	}
}

// runPar executes all of the specified functions in parallel and returns a joined error.
func runPar(ctx context.Context, r Runner, fns ...func(context.Context, Runner) error) error {
	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, fn := range fns {
		wg.Add(1)

		go func(f func(context.Context, Runner) error) {
			defer wg.Done()

			err := f(ctx, r)
			if err == nil {
				return
			}

			mu.Lock()

			errs = append(errs, err)

			mu.Unlock()
		}(fn)
	}

	wg.Wait()

	return errors.Join(errs...)
}

// ignoredDirPatterns lists directory names that file watchers should skip.
var ignoredDirPatterns = []string{
	".worktrees",
	"node_modules",
	"vendor",
	"data",
	".git",
}

// watch sets up a recursive file watcher to re-run tasks.
func watch(dir, label string, callback func(watcher.Event)) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.AddFilterHook(func(_ os.FileInfo, fullPath string) error {
		for _, pattern := range ignoredDirPatterns {
			if strings.Contains(fullPath, string(os.PathSeparator)+pattern+string(os.PathSeparator)) ||
				strings.HasSuffix(fullPath, string(os.PathSeparator)+pattern) {
				return watcher.ErrSkip
			}
		}

		return nil
	})

	go func() {
		for {
			select {
			case ev := <-w.Event:
				log.Printf("Detected change in '%s'. Running task '%s'...\n", ev.Path, label)
				callback(ev)
				log.Printf("Task '%s' completed.\n", label)
			case err := <-w.Error:
				classifyAndExit(fmtErr(err, "watcher failed"))
			case <-w.Closed:
				return
			}
		}
	}()

	classifyAndExit(fmtErr(w.AddRecursive(dir), "create watcher"))
	log.Printf("Watching '%s' for changes...\n", dir)

	go classifyAndExit(fmtErr(w.Start(100*time.Millisecond), "start watcher"))
}
