package main

import "github.com/lingrino/glen/cmd"

// version is populated at build time by goreleaser
var version = "dev"

func main() {
	cmd.Execute(version)
}
