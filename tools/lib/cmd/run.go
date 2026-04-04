package cmd

import (
	"context"
	"fmt"
	"log"
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
	runCmd.AddCommand(runWebCmd)
	rootCmd.AddCommand(runCmd)

	runAPIServerCmd.PersistentFlags().Bool("watch", false, "Watch the input directory and automatically restart the server.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run full-stack development environment",
	Run: func(cmd *cobra.Command, args []string) {
		classifyAndExit(runFullStack(cmd.Context(), runner, envProvider, args))
	},
}

func runFullStack(ctx context.Context, r Runner, ep EnvProvider, args []string) error {
	err := dockerRun(ctx, r, ep, config.ServerBinName, config.DopplerEnvName, "api-db")
	if err != nil {
		return err
	}

	offset, err := worktreePortOffset(ctx)
	if err != nil {
		return fmtErr(err, "compute port offset")
	}

	devPort := 5173 + offset
	apiPort := 5000 + offset

	devProc, err := r.Start(ctx, Cmd{
		Name: filepath.FromSlash("./node_modules/.bin/vite"),
		Args: []string{"dev", "--port", strconv.Itoa(devPort)},
		Dir:  "web-app",
	})
	if err != nil {
		return fmtErr(err, "start vite dev server")
	}

	apiProc, err := startAPIServer(ctx, r, ep, args,
		fmt.Sprintf("PUBGOLF_WEB_APP_UPSTREAM_HOST=http://localhost:%d", devPort),
		fmt.Sprintf("PUBGOLF_HOST_ORIGIN=http://localhost:%d", apiPort),
	)
	if err != nil {
		devProc.Stop()

		return err
	}

	log.Printf("Full-stack environment running (dev: :%d, API: :%d)\n", devPort, apiPort)

	go func() {
		<-shuttingDown
		apiProc.Stop()
		devProc.Stop()
	}()

	waitErr := apiProc.Wait()
	if waitErr != nil {
		log.Printf("API server exited with error: %v", waitErr)
	}

	devProc.Stop()

	return nil
}

var runAPIServerCmd = &cobra.Command{
	Use:   "api",
	Short: "Run API server",
	Run: func(cmd *cobra.Command, args []string) {
		watchableGoRun(cmd, runner, envProvider, args)
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

var runWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Run SvelteKit dev server",
	Run: func(cmd *cobra.Command, _ []string) {
		env, err := envProvider.Env(cmd.Context(), "web-app", "dev")
		classifyAndExit(fmtErr(err, "fetch web-app environment"))

		proc, startErr := runner.Start(cmd.Context(), Cmd{
			Name: filepath.FromSlash("./node_modules/.bin/vite"),
			Args: []string{"dev"},
			Dir:  "web-app",
			Env:  env,
		})
		classifyAndExit(fmtErr(startErr, "start vite dev server"))

		<-shuttingDown
		proc.Stop()
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

func watchableGoRun(cmd *cobra.Command, r Runner, ep EnvProvider, args []string) {
	watchFlag, err := cmd.Flags().GetBool("watch")
	classifyAndExit(fmtErr(err, "check '--watch' flag"))

	proc := startGoRun(cmd.Context(), r, ep, args)

	if !watchFlag {
		waitErr := proc.Wait()
		if waitErr != nil {
			log.Printf("process exited with error: %v", waitErr)
		}

		return
	}

	go func() {
		watch("api", "restart API server", func(_ watcher.Event) {
			proc.Stop()
			proc = startGoRun(context.Background(), r, ep, args)
		})
	}()

	<-shuttingDown
	proc.Stop()
}

//nolint:ireturn // Returns Process interface by design.
func startGoRun(ctx context.Context, r Runner, ep EnvProvider, args []string) Process {
	proc, err := startAPIServer(ctx, r, ep, args)
	classifyAndExit(err)

	return proc
}

//nolint:ireturn // Returns Process interface by design.
func startAPIServer(ctx context.Context, r Runner, ep EnvProvider, args []string, extraEnv ...string) (Process, error) {
	env, err := ep.Env(ctx, config.ServerBinName, config.DopplerEnvName)
	if err != nil {
		return nil, fmtErr(err, "fetch go run environment")
	}

	offset, err := worktreePortOffset(ctx)
	if err != nil {
		return nil, fmtErr(err, "compute port offset")
	}

	env = append(env,
		fmt.Sprintf("PUBGOLF_DB_PORT=%d", 5432+offset),
		fmt.Sprintf("PUBGOLF_PORT=%d", 5000+offset),
	)
	env = append(env, extraEnv...)

	bin := filepath.FromSlash("./api/cmd/" + config.ServerBinName)
	allArgs := append([]string{"run", bin}, args...)

	log.Printf("Starting '%s'...\n", config.ServerBinName)

	proc, err := r.Start(ctx, Cmd{
		Name: "go",
		Args: allArgs,
		Env:  env,
	})
	if err != nil {
		return nil, fmtErr(err, "start API server")
	}

	return proc, nil
}
