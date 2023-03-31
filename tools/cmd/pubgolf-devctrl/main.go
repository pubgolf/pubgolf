package main

import (
	"github.com/pubgolf/pubgolf/tools/lib/cmd"
)

var toolsHash string

func main() {
	cmd.Execute(toolsHash, cmd.CLIConfig{
		ProjectName:  "pubgolf",
		DBDriver:     cmd.PostgreSQL,
		EnvVarPrefix: "PUBGOLF_",
	})
}
