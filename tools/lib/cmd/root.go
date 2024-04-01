// Package cmd contains handlers for each CLI command.
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	"golang.org/x/mod/sumdb/dirhash"
)

// Internal plumbing.
var (
	// beginShutdown indicates we've received an OS signal to begin shutting down.
	beginShutdown = make(chan os.Signal, 1)
	// shuttingDown is a broadcast channel that closes to tell all processes to begin cleanup.
	shuttingDown = make(chan struct{})
)

// Package parameters.
var (
	// installedToolsHash holds the hash of the source code for the currently installed version of the binary.
	installedToolsHash string
	// config holds the provided configuration options.
	config CLIConfig
)

func init() {
	signal.Notify(beginShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

// Execute is the entrypoint for calling the CLI.
func Execute(toolsDirHash string, c CLIConfig) {
	installedToolsHash = toolsDirHash

	c.setDefaults()
	config = c

	log.SetPrefix(fmt.Sprintf("[%s] ", config.CLIName))
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   config.CLIName,
	Short: "DevCtrl is a task runner for local dev",
	Long:  `An opinionated task runner for personal use by @thedeerchild`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip the update warning for the update command itself
		if cmd.CommandPath() != " update" {
			checkVersion()
		}

		go func() {
			<-beginShutdown
			close(shuttingDown)
		}()
	},
}

func checkVersion() {
	curToolsHash, err := dirhash.HashDir("tools", "", dirhash.DefaultHash)
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
		log.Panicf("%s: %v", msg, err.Error())
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
