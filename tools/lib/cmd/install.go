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
	Run: func(cmd *cobra.Command, _ []string) {
		fns := []func(context.Context, Runner) error{
			func(ctx context.Context, r Runner) error {
				err := installWithHomebrew(ctx, r, "gnupg")
				if err != nil {
					return err
				}

				return installWithHomebrew(ctx, r, "dopplerhq/cli/doppler")
			},
			func(ctx context.Context, r Runner) error {
				return installWithHomebrew(ctx, r, "bufbuild/buf/buf")
			},
		}

		for _, pkg := range goTools {
			fns = append(fns, func(ctx context.Context, r Runner) error {
				return installWithGolang(ctx, r, pkg)
			})
		}

		classifyAndExit(runPar(cmd.Context(), runner, fns...))
	},
}

func generateGoInstallers() []*cobra.Command {
	cs := make([]*cobra.Command, 0, len(goTools))

	for n, v := range goTools {
		cs = append(cs, &cobra.Command{
			Use:   n,
			Short: fmt.Sprintf("Install the %q CLI tool", n),
			Run: func(cmd *cobra.Command, _ []string) {
				classifyAndExit(installWithGolang(cmd.Context(), runner, v))
			},
		})
	}

	return cs
}

var installDopplerCmd = &cobra.Command{
	Use:   "doppler",
	Short: "Install the `doppler` CLI tool",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(installWithHomebrew(cmd.Context(), runner, "gnupg"))
		classifyAndExit(installWithHomebrew(cmd.Context(), runner, "dopplerhq/cli/doppler"))
	},
}

var installBufCmd = &cobra.Command{
	Use:   "buf",
	Short: "Install the `buf` CLI tool",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(installWithHomebrew(cmd.Context(), runner, "bufbuild/buf/buf"))
	},
}

func installWithGolang(ctx context.Context, r Runner, pkg string) error {
	err := r.Run(ctx, Cmd{
		Name: "go",
		Args: []string{"install", pkg},
	})
	if err != nil {
		return fmtErr(err, fmt.Sprintf("run go install cmd for package '%s'", pkg))
	}

	return nil
}

func installWithHomebrew(ctx context.Context, r Runner, pkg string) error {
	err := r.Run(ctx, Cmd{
		Name: "brew",
		Args: []string{"install", pkg},
	})
	if err != nil {
		return fmtErr(err, fmt.Sprintf("run brew install cmd for package '%s'", pkg))
	}

	return nil
}
