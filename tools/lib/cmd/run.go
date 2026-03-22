package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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
	runCmd.PersistentFlags().Int("port-offset", 0, "Override the port offset for this worktree (sets PUBGOLF_PORT_OFFSET)")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all server executables",
	Run: func(cmd *cobra.Command, args []string) {
		applyPortOffsetFlag(cmd)
		classifyAndExit(dockerRun(cmd.Context(), runner, envProvider, config.ServerBinName, config.DopplerEnvName,
			"api-db",
		))

		binPath := filepath.FromSlash("./api/cmd/" + config.ServerBinName)
		watchableGoRun(cmd, runner, envProvider, config.ServerBinName, config.DopplerEnvName, binPath, args)
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
		classifyAndExit(dockerRun(cmd.Context(), runner, envProvider, config.ServerBinName, config.DopplerEnvName,
			"api-db",
			// "api-blob-storage",
		))
	},
}

var runAPIDatabaseCmd = &cobra.Command{
	Use:   "api-db",
	Short: "Run API server's DB instance",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(dockerRun(cmd.Context(), runner, envProvider, config.ServerBinName, config.DopplerEnvName, "api-db"))
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

	offset, offsetErr := worktreePortOffset(ctx)
	if offsetErr != nil {
		return fmtErr(offsetErr, "compute port offset")
	}

	env = append(env,
		fmt.Sprintf("PUBGOLF_DB_PORT=%d", 5432+offset),
		fmt.Sprintf("PUBGOLF_PORT=%d", 5000+offset),
		"PUBGOLF_DB_HOST_DATA_PATH="+filepath.Join(projectRoot, worktreeDataDir(ctx, "data/postgres")),
	)

	projectName := worktreeDockerProject(ctx)

	args := []string{
		"compose",
		"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
		"--project-name", projectName,
		"up",
		"--detach",
		"--",
	}
	args = append(args, services...)

	runErr := r.Run(ctx, Cmd{
		Name: "docker",
		Args: args,
		Env:  env,
	})
	if runErr != nil {
		return fmtErr(runErr, "run docker-compose up cmd")
	}

	// Post-start reminder for worktree users.
	if slug, _ := worktreeSlug(ctx); slug != "" {
		log.Printf("Started services for worktree %q (DB: port %d).\n"+
			"  Run 'pubgolf-devctrl stop' before removing this worktree.", slug, 5432+offset)
	}

	return nil
}

func watchableGoRun(cmd *cobra.Command, r Runner, ep EnvProvider, project, envCfg, bin string, args []string) {
	watchFlag, err := cmd.Flags().GetBool("watch")
	classifyAndExit(fmtErr(err, "check '--watch' flag"))

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

//nolint:ireturn // Returns Process interface by design.
func startGoRun(ctx context.Context, r Runner, ep EnvProvider, project, envCfg, bin string, args []string) Process {
	env, err := ep.Env(ctx, project, envCfg)
	classifyAndExit(fmtErr(err, "fetch go run environment"))

	offset, offsetErr := worktreePortOffset(ctx)
	classifyAndExit(fmtErr(offsetErr, "compute port offset"))

	env = append(env,
		fmt.Sprintf("PUBGOLF_DB_PORT=%d", 5432+offset),
		fmt.Sprintf("PUBGOLF_PORT=%d", 5000+offset),
	)

	allArgs := append([]string{"run", bin}, args...)

	log.Printf("Starting '%s'...\n", bin)

	proc, startErr := r.Start(ctx, Cmd{
		Name: "go",
		Args: allArgs,
		Env:  env,
	})
	classifyAndExit(fmtErr(startErr, "execute `go run ...` command"))

	return proc
}

// applyPortOffsetFlag reads the --port-offset flag and sets PUBGOLF_PORT_OFFSET
// in the process environment so worktreePortOffset() picks it up.
func applyPortOffsetFlag(cmd *cobra.Command) {
	offset, err := cmd.Flags().GetInt("port-offset")
	classifyAndExit(fmtErr(err, "check '--port-offset' flag"))

	if offset > 0 {
		classifyAndExit(fmtErr(os.Setenv("PUBGOLF_PORT_OFFSET", strconv.Itoa(offset)), "set PUBGOLF_PORT_OFFSET"))
	}
}
