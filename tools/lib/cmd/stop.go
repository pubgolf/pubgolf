package cmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all background processes started with `devctrl run ...`",
	Run: func(cmd *cobra.Command, _ []string) {
		dopplerDockerStop(cmd.Context(), runner, config.ServerBinName, config.DopplerEnvName)
	},
}

func dopplerDockerStop(ctx context.Context, r Runner, project, env string) {
	guard(r.Run(ctx, Cmd{
		Name: "doppler",
		Args: []string{
			"run",
			"--project", project,
			"--config", env,
			"--",
			"docker-compose",
			"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
			"down",
		},
	}), "execute `docker-compose down ...` command")
}
