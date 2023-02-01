package main

import (
	"github.com/pubgolf/pubgolf/tools/lib/cmd"
)

var toolsHash string

func main() {
	projectName := "pubgolf"
	cmd.Execute(projectName, toolsHash)
}
