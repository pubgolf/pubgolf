package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var goTools = map[string]string{
	"migrate":       "github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1",
	"sqlc":          "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0",
	"mockery":       "github.com/vektra/mockery/v2@v2.53.6",
	"ifacemaker":    "github.com/vburenin/ifacemaker@v1.3.0",
	"enumer":        "github.com/dmarkham/enumer@v1.6.3",
	"golangci-lint": "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4",
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
			Run: func(cmd *cobra.Command, _ []string) {
				installWithGolang(cmd.Context(), runner, v)
			},
		})
	}

	return cs
}

var installDopplerCmd = &cobra.Command{
	Use:   "doppler",
	Short: "Install the `doppler` CLI tool",
	Run: func(cmd *cobra.Command, _ []string) {
		installWithHomebrew(cmd.Context(), runner, "gnupg")
		installWithHomebrew(cmd.Context(), runner, "dopplerhq/cli/doppler")
	},
}

var installBufCmd = &cobra.Command{
	Use:   "buf",
	Short: "Install the `buf` CLI tool",
	Run: func(cmd *cobra.Command, _ []string) {
		installWithHomebrew(cmd.Context(), runner, "bufbuild/buf/buf")
	},
}

func installWithGolang(ctx context.Context, r Runner, pkg string) {
	guard(r.Run(ctx, Cmd{
		Name: "go",
		Args: []string{"install", pkg},
	}), fmt.Sprintf("execute `go install ...` command for package '%s'", pkg))
}

func installWithHomebrew(ctx context.Context, r Runner, pkg string) {
	guard(r.Run(ctx, Cmd{
		Name: "brew",
		Args: []string{"install", pkg},
	}), fmt.Sprintf("execute `brew install ...` command for package '%s'", pkg))
}
