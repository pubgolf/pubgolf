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
		guard(dopplerDockerStop(cmd.Context(), runner, config.ServerBinName, config.DopplerEnvName),
			"execute `docker-compose down ...` command")
	},
}

func dopplerDockerStop(ctx context.Context, r Runner, project, env string) error {
	err := r.Run(ctx, Cmd{
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
	})
	if err != nil {
		return fmtErr(err, "run docker-compose down cmd")
	}

	return nil
}
