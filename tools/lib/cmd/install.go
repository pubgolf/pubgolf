package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var goTools = map[string]string{
	"migrate":    "github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0",
	"sqlc":       "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0",
	"mockery":    "github.com/vektra/mockery/v2@v2.38.0",
	"ifacemaker": "github.com/vburenin/ifacemaker@v1.2.1",
	"enumer":     "github.com/dmarkham/enumer@v1.5.9",
}

func init() {
	installCmd.AddCommand(installDopplerCmd)
	installCmd.AddCommand(installBufCmd)

	for _, cmd := range generateGoInstallers() {
		installCmd.AddCommand(cmd)
	}

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run all tool installation sub-tasks",
	Run: func(cmd *cobra.Command, args []string) {
		installers := []*cobra.Command{
			installDopplerCmd,
			installBufCmd,
		}

		installers = append(installers, generateGoInstallers()...)

		runPar(cmd, args, installers...)
	},
}

func generateGoInstallers() []*cobra.Command {
	cs := make([]*cobra.Command, 0, len(goTools))

	for n, v := range goTools {
		cs = append(cs, &cobra.Command{
			Use:   n,
			Short: fmt.Sprintf("Install the %q CLI tool", n),
			Run: func(_ *cobra.Command, _ []string) {
				installWithGolang(v)
			},
		})
	}

	return cs
}

var installDopplerCmd = &cobra.Command{
	Use:   "doppler",
	Short: "Install the `doppler` CLI tool",
	Run: func(_ *cobra.Command, _ []string) {
		installWithHomebrew("gnupg")
		installWithHomebrew("dopplerhq/cli/doppler")
	},
}

var installBufCmd = &cobra.Command{
	Use:   "buf",
	Short: "Install the `buf` CLI tool",
	Run: func(_ *cobra.Command, _ []string) {
		installWithHomebrew("bufbuild/buf/buf")
	},
}

func installWithGolang(pkg string) {
	installer := exec.Command("go", "install", pkg)
	installer.Stdout = os.Stdout
	installer.Stderr = os.Stderr
	guard(installer.Run(), fmt.Sprintf("execute `go install ...` command for package '%s'", pkg))
}

func installWithHomebrew(pkg string) {
	installer := exec.Command("brew", "install", pkg)
	installer.Stdout = os.Stdout
	installer.Stderr = os.Stderr
	guard(installer.Run(), fmt.Sprintf("execute `brew install ...` command for package '%s'", pkg))
}
