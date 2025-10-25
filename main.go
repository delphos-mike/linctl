package main

import (
	_ "embed"

	"github.com/delphos-mike/linctl/cmd"
)

//go:embed README.md
var readmeContents string

func main() {
	cmd.SetReadmeContents(readmeContents)
	cmd.Execute()
}
