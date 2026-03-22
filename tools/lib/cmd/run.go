package cmd

import (
	"context"
	"log"
	"path/filepath"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.AddCommand(runAPIServerCmd)
	runCmd.AddCommand(runAPIBgCmd)
	runCmd.AddCommand(runAPIDatabaseCmd)
	// runCmd.AddCommand(runAPIMinioCmd)
	rootCmd.AddCommand(runCmd)

	runAPIServerCmd.PersistentFlags().Bool("watch", false, "Watch the input directory and automatically restart the server.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all server executables",
	Run: func(cmd *cobra.Command, args []string) {
		runAPIBgCmd.Run(cmd, args)
		runPar(cmd, args,
			runAPIServerCmd,
		)
	},
}

var runAPIServerCmd = &cobra.Command{
	Use:   "api",
	Short: "Run API server",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := filepath.FromSlash("./api/cmd/" + config.ServerBinName)
		watchableDopplerGoRun(cmd, runner, config.ServerBinName, config.DopplerEnvName, binPath, args)
	},
}

var runAPIBgCmd = &cobra.Command{
	Use:   "bg",
	Short: "Run all supporting services for the API server",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(dopplerDockerRun(cmd.Context(), runner, config.ServerBinName, config.DopplerEnvName,
			"api-db",
			// "api-blob-storage",
		), "execute `docker-compose up ...` command")
	},
}

var runAPIDatabaseCmd = &cobra.Command{
	Use:   "api-db",
	Short: "Run API server's DB instance",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(dopplerDockerRun(cmd.Context(), runner, config.ServerBinName, config.DopplerEnvName, "api-db"),
			"execute `docker-compose up ...` command")
	},
}

// var runAPIMinioCmd = &cobra.Command{
// 	Use:   "api-minio",
// 	Short: "Run API server's blob storage (Minio) instance",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dopplerDockerRun(runner, config.ServerBinName, config.DopplerEnvName, "api-blob-storage")
// 	},
// }

func dopplerDockerRun(ctx context.Context, r Runner, project, env string, services ...string) error {
	args := []string{
		"run",
		"--project", project,
		"--config", env,
		"--",
		"docker-compose",
		"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
		"up",
		"--detach",
		"--",
	}
	args = append(args, services...)

	err := r.Run(ctx, Cmd{
		Name: "doppler",
		Args: args,
	})
	if err != nil {
		return fmtErr(err, "run docker-compose up cmd")
	}

	return nil
}

func watchableDopplerGoRun(cmd *cobra.Command, r Runner, project, env, bin string, args []string) {
	watchFlag, err := cmd.Flags().GetBool("watch")
	guard(err, "check '--watch' flag")

	// Start initial process
	proc := dopplerGoRun(cmd.Context(), r, project, env, bin, args, !watchFlag)

	// Launch watcher, if applicable.
	if watchFlag {
		go func() {
			watch("api", "restart API server", func(_ watcher.Event) {
				// Stop the old process and keep track of the new one.
				proc.Stop()
				proc = dopplerGoRun(context.Background(), r, project, env, bin, args, false)
			})
		}()
	}

	// Hold process open
	<-shuttingDown
	proc.Stop()
}

func dopplerGoRun(ctx context.Context, r Runner, project, env, bin string, args []string, stopOnExit bool) Process { //nolint:ireturn // Returns Process interface by design.
	allArgs := append(
		[]string{
			"run",
			"--project", project,
			"--config", env,
			"--",
			"go", "run", bin,
		}, args...)

	log.Printf("Starting '%s'...\n", bin)

	proc, err := r.Start(ctx, Cmd{
		Name: "doppler",
		Args: allArgs,
	})
	guard(err, "execute `go run ...` command")

	if stopOnExit {
		go func() {
			waitErr := proc.Wait()
			if waitErr != nil {
				log.Printf("process exited with error: %v", waitErr)
			}

			triggerShutdown()
		}()
	}

	return proc
}
