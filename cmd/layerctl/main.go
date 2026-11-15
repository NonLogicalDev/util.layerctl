// Modified from abg-OAI/codex: module import path changed for this extraction.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/NonLogicalDev/util.layerctl/internal/cli"
)

// Set by release and package builds; local builds remain identifiable.
var version = "0.0.0-dev"
var commit = "unknown"

// Go records VCS metadata for ordinary local builds when the checkout is available.
func buildCommit() string {
	if commit != "unknown" {
		return commit
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return commit
	}
	revision, modified := commit, false
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	if revision != "unknown" && modified {
		return revision + "-dirty"
	}
	return revision
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("layerctl %s (%s)\n", version, buildCommit())
		return
	}
	logger := log.New(os.Stderr, "layerctl: ", 0)
	if err := cli.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		logger.Print(err)
		os.Exit(1)
	}
}
