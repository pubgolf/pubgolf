package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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
		watchableDopplerGoRun(cmd, config.ServerBinName, config.DopplerEnvName, binPath, args)
	},
}

var runAPIBgCmd = &cobra.Command{
	Use:   "bg",
	Short: "Run all supporting services for the API server",
	Run: func(cmd *cobra.Command, args []string) {
		dopplerDockerRun(config.ServerBinName, config.DopplerEnvName,
			"api-db",
			// "api-blob-storage",
		)
	},
}

var runAPIDatabaseCmd = &cobra.Command{
	Use:   "api-db",
	Short: "Run API server's DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		dopplerDockerRun(config.ServerBinName, config.DopplerEnvName, "api-db")
	},
}

// var runAPIMinioCmd = &cobra.Command{
// 	Use:   "api-minio",
// 	Short: "Run API server's blob storage (Minio) instance",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		dopplerDockerRun(config.ServerBinName, config.DopplerEnvName, "api-blob-storage")
// 	},
// }

func dopplerDockerRun(project, env string, services ...string) {
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

	doppler := exec.Command("doppler", args...)

	doppler.Stdout = os.Stdout
	doppler.Stderr = os.Stderr
	doppler.Stdin = os.Stdin

	guard(doppler.Run(), "execute `docker-compose up ...` command")
}

func watchableDopplerGoRun(cmd *cobra.Command, project, env, bin string, args []string) {
	watchFlag, err := cmd.Flags().GetBool("watch")
	guard(err, "check '--watch' flag")

	// Start initial process
	stopFn := dopplerGoRun(project, env, bin, args, !watchFlag)

	// Launch watcher, if applicable.
	if watchFlag {
		go func() {
			watch("api", "restart API server", func(ev watcher.Event) {
				// Start the old process and keep track of the new cleanup handler.
				stopFn()
				stopFn = dopplerGoRun(project, env, bin, args, false)
			})
		}()
	}

	// Hold process open
	<-shuttingDown
	stopFn()
}

func dopplerGoRun(project, env, bin string, args []string, stopOnExit bool) func() {
	allArgs := append(
		[]string{
			"run",
			"--project", project,
			"--config", env,
			"--",
			"go", "run", bin,
		}, args...)

	doppler := exec.Command("doppler", allArgs...)

	doppler.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	doppler.Stdout = os.Stdout
	doppler.Stderr = os.Stderr
	doppler.Stdin = os.Stdin

	log.Printf("Starting '%s'...\n", bin)
	guard(doppler.Start(), "execute `go run ...` command")

	if stopOnExit {
		go func() {
			doppler.Wait()
			close(beginShutdown)
		}()
	}

	return func() {
		log.Printf("Calling Doppler shutdown for pid %d...", doppler.Process.Pid)
		pgid, err := syscall.Getpgid(doppler.Process.Pid)
		if err != nil {
			log.Println("Could not get pgid. Assuming process has already shut down.")
			return
		}
		guard(syscall.Kill(-pgid, syscall.SIGKILL), "send SIGINT to running process")
	}
}
