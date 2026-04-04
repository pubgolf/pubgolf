package cmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	buildCmd.AddCommand(buildWebCmd)
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build all targets",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(buildWeb(cmd.Context(), runner))
	},
}

var buildWebCmd = &cobra.Command{
	Use:   "web",
	Short: "Build the web app with Vite",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(buildWeb(cmd.Context(), runner))
	},
}

func buildWeb(ctx context.Context, r Runner, extraEnv ...string) error {
	err := r.Run(ctx, Cmd{
		Name: filepath.FromSlash("./node_modules/.bin/vite"),
		Args: []string{"build"},
		Dir:  "web-app",
		Env:  extraEnv,
	})
	if err != nil {
		return fmtErr(err, "build web-app")
	}

	return nil
}
