// Package cmd contains handlers for each CLI command.
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
}

func checkVersion() {
	curToolsHash, err := dirhash.HashDir(filepath.Join(projectRoot, "tools"), "", dirhash.DefaultHash)
	guard(err, "hash tools dir")

	if installedToolsHash != curToolsHash {
		log.Printf(`The installed version of %[1]s is out of date. Run the following to update:

%[1]s update

`, config.CLIName)
	}
}

// guard logs and exits on error.
func guard(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err.Error())
	}
}

// runPar executes all of the specified commands in parallel.
func runPar(cmd *cobra.Command, args []string, commands ...*cobra.Command) {
	var wg sync.WaitGroup

	wg.Add(len(commands))

	for _, c := range commands {
		go func(cmd *cobra.Command, args []string, c *cobra.Command) {
			defer wg.Done()

			log.Printf("Running `%s`...\n", c.Use)
			c.Run(cmd, args)
		}(cmd, args, c)
	}

	wg.Wait()
}

// watch sets up a recursive file watcher to re-run tasks.
func watch(dir, label string, callback func(watcher.Event)) {
	w := watcher.New()
	w.SetMaxEvents(1)

	go func() {
		for {
			select {
			case ev := <-w.Event:
				log.Printf("Detected change in '%s'. Running task '%s'...\n", ev.Path, label)
				callback(ev)
				log.Printf("Task '%s' completed.\n", label)
			case err := <-w.Error:
				guard(err, "watcher failed")
			case <-w.Closed:
				return
			}
		}
	}()

	guard(w.AddRecursive(dir), "create watcher")
	log.Printf("Watching '%s' for changes...\n", dir)

	go guard(w.Start(100*time.Millisecond), "start watcher")
}
