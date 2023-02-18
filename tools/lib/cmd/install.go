package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	installCmd.AddCommand(installDopplerCmd)
	installCmd.AddCommand(installMigrateCmd)
	installCmd.AddCommand(installSQLcCmd)
	installCmd.AddCommand(installMockeryCmd)
	installCmd.AddCommand(installIfacemakerCmd)
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Run all tool installation sub-tasks",
	Run: func(cmd *cobra.Command, args []string) {
		installers := []*cobra.Command{
			installDopplerCmd,
			installMigrateCmd,
			installSQLcCmd,
			installMockeryCmd,
			installIfacemakerCmd,
		}

		runPar(cmd, args, installers...)
	},
}

var installDopplerCmd = &cobra.Command{
	Use:   "doppler",
	Short: "Install the `doppler` CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		installWithHomebrew("libpq")
		installer := exec.Command("bash", "-c", "(curl -Ls --tlsv1.2 --proto '=https' --retry 3 https://cli.doppler.com/install.sh || wget -t 3 -qO- https://cli.doppler.com/install.sh) | sh -s -- --verify-signature")
		installer.Stdout = os.Stdout
		installer.Stderr = os.Stderr
		guard(installer.Run(), "execute cURL installation command for Doppler CLI")
	},
}

var installMigrateCmd = &cobra.Command{
	Use:   "golang-migrate",
	Short: "Install the `migrate` CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		installWithGolang("github.com/golang-migrate/migrate/v4/cmd/migrate")
	},
}

var installSQLcCmd = &cobra.Command{
	Use:   "sqlc",
	Short: "Install the `sqlc` CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		installWithGolang("github.com/kyleconroy/sqlc/cmd/sqlc")
	},
}

var installMockeryCmd = &cobra.Command{
	Use:   "mockery",
	Short: "Install the `mockery` CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		installWithGolang("github.com/vektra/mockery/v2")
	},
}

var installIfacemakerCmd = &cobra.Command{
	Use:   "ifacemaker",
	Short: "Install the `ifacemaker` CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		installWithGolang("github.com/vburenin/ifacemaker")
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
