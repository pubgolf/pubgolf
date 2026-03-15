package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: fmt.Sprintf("Stop all background processes started with `%s run ...`", config.CLIName),
	Run: func(cmd *cobra.Command, _ []string) {
		dopplerDockerStop(cmd.Context(), config.ServerBinName, config.DopplerEnvName)
	},
}

func dopplerDockerStop(ctx context.Context, project, env string) {
	doppler := exec.CommandContext(ctx, "doppler",
		"run",
		"--project", project,
		"--config", env,
		"--",
		"docker-compose",
		"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
		"down")

	doppler.Stdout = os.Stdout
	doppler.Stderr = os.Stderr
	doppler.Stdin = os.Stdin

	guard(doppler.Run(), "execute `docker-compose up ...` command")
}
