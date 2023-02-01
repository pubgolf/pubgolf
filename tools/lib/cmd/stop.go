package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: fmt.Sprintf("Stop all background processes started with `%s run ...`", commandName),
	Run: func(cmd *cobra.Command, args []string) {
		dopplerDockerStop(serverBinName, "dev")
	},
}

func dopplerDockerStop(project, env string, services ...string) {
	doppler := exec.Command("doppler",
		"run",
		"--project", project,
		"--config", env,
		"--",
		"docker-compose",
		"--file", "docker-compose.dev.yaml",
		"down")

	doppler.Stdout = os.Stdout
	doppler.Stderr = os.Stderr
	doppler.Stdin = os.Stdin

	guard(doppler.Run(), "execute `docker-compose up ...` command")
}
