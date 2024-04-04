package cmd

import (
	"os"
	"os/exec"
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
	Run: func(_ *cobra.Command, _ []string) {
		curToolsHash, err := dirhash.HashDir("tools", "", dirhash.DefaultHash)
		guard(err, "hash tools dir")

		// TODO: include go.sum in generated hash.

		compiler := exec.Command("go", //nolint:gosec // Input is not dynamically provided by end-user.
			"install",
			"-ldflags", "-X main.toolsHash="+curToolsHash,
			filepath.FromSlash("./tools/cmd/"+config.CLIName),
		)
		compiler.Stderr = os.Stderr
		compiler.Stdout = os.Stdout

		guard(compiler.Run(), "execute `go install ...` command")
	},
}
