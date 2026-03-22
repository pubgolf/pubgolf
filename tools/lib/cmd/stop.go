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
		guard(dockerStop(cmd.Context(), runner),
			"execute `docker-compose down ...` command")
	},
}

func dockerStop(ctx context.Context, r Runner) error {
	// docker-compose down doesn't need secrets; run without env injection.
	err := r.Run(ctx, Cmd{
		Name: "docker-compose",
		Args: []string{
			"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
			"down",
		},
	})
	if err != nil {
		return fmtErr(err, "run docker-compose down cmd")
	}

	return nil
}
