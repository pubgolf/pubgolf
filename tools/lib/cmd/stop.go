package cmd

import (
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
	Run: func(_ *cobra.Command, _ []string) {
		dopplerDockerStop(config.ServerBinName, config.DopplerEnvName)
	},
}

func dopplerDockerStop(project, env string) {
	doppler := exec.Command("doppler", //nolint:gosec // Input is not dynamically provided by end-user.
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
