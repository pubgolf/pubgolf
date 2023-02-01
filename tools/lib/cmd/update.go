package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		curToolsHash, err := dirhash.HashDir("tools", "", dirhash.DefaultHash)
		guard(err, "hash tools dir")

		compiler := exec.Command("go",
			"install",
			"-ldflags", fmt.Sprintf("-X main.toolsHash=%s", curToolsHash),
			filepath.FromSlash("./tools/cmd/"+commandName),
		)
		compiler.Stderr = os.Stderr
		compiler.Stdout = os.Stdout

		guard(compiler.Run(), "execute `go install ...` command")
	},
}
