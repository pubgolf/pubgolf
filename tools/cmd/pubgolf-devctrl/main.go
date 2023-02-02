package main

import (
	"github.com/pubgolf/pubgolf/tools/lib/cmd"
)

var toolsHash string

func main() {
	cmd.Execute(toolsHash, cmd.DevCtrlConfig{
		ProjectName: "pubgolf",
		DBDriver:    cmd.SQLite3,
	})
}
