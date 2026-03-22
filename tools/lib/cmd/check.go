package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	checkCmd.AddCommand(checkGoCmd)
	checkCmd.AddCommand(checkWebCmd)
	checkCmd.AddCommand(checkProtoCmd)
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run all linting and static analysis sub-tasks",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(runPar(cmd.Context(), runner, checkGo, checkWeb, checkProto))
	},
}

var checkGoCmd = &cobra.Command{
	Use:   "go",
	Short: "Run golangci-lint on Go code",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(checkGo(cmd.Context(), runner))
	},
}

func checkGo(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: "golangci-lint",
		Args: []string{"run", "./api/...", "./tools/..."},
	})
	if err != nil {
		return fmtErr(err, "run golangci-lint cmd")
	}

	return nil
}

var checkWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Run ESLint, Prettier, and svelte-check on web code",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(checkWeb(cmd.Context(), runner))
	},
}

func checkWeb(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: "npm",
		Args: []string{"run", "ci:lint"},
		Dir:  "web-app",
	})
	if err != nil {
		return fmtErr(err, "run npm ci:lint cmd")
	}

	err = r.Run(ctx, Cmd{
		Name: "npm",
		Args: []string{"run", "check"},
		Dir:  "web-app",
	})
	if err != nil {
		return fmtErr(err, "run npm check cmd")
	}

	return nil
}

var checkProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Run buf lint and format check on proto files",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(checkProto(cmd.Context(), runner))
	},
}

func checkProto(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: "buf",
		Args: []string{"lint"},
	})
	if err != nil {
		return fmtErr(err, "run buf lint cmd")
	}

	err = r.Run(ctx, Cmd{
		Name: "buf",
		Args: []string{"format", "./proto/", "--diff", "--exit-code"},
	})
	if err != nil {
		return fmtErr(err, "run buf format cmd")
	}

	return nil
}
