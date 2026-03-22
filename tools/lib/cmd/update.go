package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/mod/sumdb/dirhash"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Compiles the latest changes into the devctl binary",
	Run: func(cmd *cobra.Command, _ []string) {
		curToolsHash, err := dirhash.HashDir("tools", "", dirhash.DefaultHash)
		classifyAndExit(fmtErr(err, "hash tools dir"))

		// TODO: include go.sum in generated hash.

		classifyAndExit(fmtErr(runner.Run(cmd.Context(), Cmd{
			Name: "go",
			Args: []string{
				"install",
				"-ldflags", "-X main.toolsHash=" + curToolsHash,
				filepath.FromSlash("./tools/cmd/" + config.CLIName),
			},
		}), "execute `go install ...` command"))
	},
}
