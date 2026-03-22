package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

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
	Run: func(cmd *cobra.Command, args []string) {
		runPar(cmd, args, checkGoCmd, checkWebCmd, checkProtoCmd)
	},
}

var checkGoCmd = &cobra.Command{
	Use:   "go",
	Short: "Run golangci-lint on Go code",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(checkGo(cmd.Context()), "execute `golangci-lint run` command")
	},
}

func checkGo(ctx context.Context) error {
	lint := exec.CommandContext(ctx, "golangci-lint", "run", "--exclude-dirs", `\.worktrees`, "./api/...", "./tools/...")
	lint.Stdout = os.Stdout
	lint.Stderr = os.Stderr

	err := lint.Run()
	if err != nil {
		return fmt.Errorf("run golangci-lint cmd: %w", err)
	}

	return nil
}

var checkWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Run ESLint, Prettier, and svelte-check on web code",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(checkWeb(cmd.Context()), "execute web lint commands")
	},
}

func checkWeb(ctx context.Context) error {
	ciLint := exec.CommandContext(ctx, "npm", "run", "ci:lint")
	ciLint.Dir = "web-app"
	ciLint.Stdout = os.Stdout
	ciLint.Stderr = os.Stderr

	err := ciLint.Run()
	if err != nil {
		return fmt.Errorf("run npm ci:lint cmd: %w", err)
	}

	check := exec.CommandContext(ctx, "npm", "run", "check")
	check.Dir = "web-app"
	check.Stdout = os.Stdout
	check.Stderr = os.Stderr

	err = check.Run()
	if err != nil {
		return fmt.Errorf("run npm check cmd: %w", err)
	}

	return nil
}

var checkProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Run buf lint and format check on proto files",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(checkProto(cmd.Context()), "execute proto lint commands")
	},
}

func checkProto(ctx context.Context) error {
	lint := exec.CommandContext(ctx, "buf", "lint")
	lint.Stdout = os.Stdout
	lint.Stderr = os.Stderr

	err := lint.Run()
	if err != nil {
		return fmt.Errorf("run buf lint cmd: %w", err)
	}

	format := exec.CommandContext(ctx, "buf", "format", "./proto/", "--diff", "--exit-code")
	format.Stdout = os.Stdout
	format.Stderr = os.Stderr

	err = format.Run()
	if err != nil {
		return fmt.Errorf("run buf format cmd: %w", err)
	}

	return nil
}
