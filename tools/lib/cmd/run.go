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
		watchableGoRun(cmd, runner, envProvider, config.ServerBinName, config.DopplerEnvName, binPath, args)
	},
}

var runAPIBgCmd = &cobra.Command{
	Use:   "bg",
	Short: "Run all supporting services for the API server",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(dockerRun(cmd.Context(), runner, envProvider, config.ServerBinName, config.DopplerEnvName,
			"api-db",
			// "api-blob-storage",
		), "execute `docker-compose up ...` command")
	},
}

var runAPIDatabaseCmd = &cobra.Command{
	Use:   "api-db",
	Short: "Run API server's DB instance",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(dockerRun(cmd.Context(), runner, envProvider, config.ServerBinName, config.DopplerEnvName, "api-db"),
			"execute `docker-compose up ...` command")
	},
}

// var runAPIMinioCmd = &cobra.Command{
// 	Use:   "api-minio",
// 	Short: "Run API server's blob storage (Minio) instance",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dockerRun(runner, envProvider, config.ServerBinName, config.DopplerEnvName, "api-blob-storage")
// 	},
// }

func dockerRun(ctx context.Context, r Runner, ep EnvProvider, project, envCfg string, services ...string) error {
	env, err := ep.Env(ctx, project, envCfg)
	if err != nil {
		return fmtErr(err, "fetch docker-compose environment")
	}

	args := []string{
		"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
		"up",
		"--detach",
		"--",
	}
	args = append(args, services...)

	runErr := r.Run(ctx, Cmd{
		Name: "docker-compose",
		Args: args,
		Env:  env,
	})
	if runErr != nil {
		return fmtErr(runErr, "run docker-compose up cmd")
	}

	return nil
}

func watchableGoRun(cmd *cobra.Command, r Runner, ep EnvProvider, project, envCfg, bin string, args []string) {
	watchFlag, err := cmd.Flags().GetBool("watch")
	guard(err, "check '--watch' flag")

	// Start initial process
	proc := startGoRun(cmd.Context(), r, ep, project, envCfg, bin, args)

	if !watchFlag {
		// Wait for process to exit, then shut down.
		waitErr := proc.Wait()
		if waitErr != nil {
			log.Printf("process exited with error: %v", waitErr)
		}

		return
	}

	// Launch watcher to restart the process on file changes.
	go func() {
		watch("api", "restart API server", func(_ watcher.Event) {
			proc.Stop()
			proc = startGoRun(context.Background(), r, ep, project, envCfg, bin, args)
		})
	}()

	// Hold process open until shutdown signal.
	<-shuttingDown
	proc.Stop()
}

func startGoRun(ctx context.Context, r Runner, ep EnvProvider, project, envCfg, bin string, args []string) Process { //nolint:ireturn // Returns Process interface by design.
	env, err := ep.Env(ctx, project, envCfg)
	guard(err, "fetch go run environment")

	allArgs := append([]string{"run", bin}, args...)

	log.Printf("Starting '%s'...\n", bin)

	proc, startErr := r.Start(ctx, Cmd{
		Name: "go",
		Args: allArgs,
		Env:  env,
	})
	guard(startErr, "execute `go run ...` command")

	return proc
}
