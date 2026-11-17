// Modified from abg-OAI/codex: module import path changed for this extraction.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/NonLogicalDev/util.layerctl/internal/cli"
)

// Set by release and package builds; local builds remain identifiable.
var version = "0.0.0-dev"
var commit = "unknown"

func main() {
	logger := log.New(os.Stderr, "layerctl: ", 0)
	if err := cli.Run(context.Background(), os.Args[1:], logger, os.Getwd, versionInfo, kong.Exit(os.Exit)); err != nil {
		logger.Print(err)
		os.Exit(1)
	}
}

func versionInfo() string {
	info, _ := debug.ReadBuildInfo()
	return fmt.Sprintf("%s (%s)", buildVersion(version, info), buildCommit(info))
}

// A versioned go install records its module version without release linker flags.
func buildVersion(embedded string, info *debug.BuildInfo) string {
	if embedded != "0.0.0-dev" {
		return embedded
	}
	if info != nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return embedded
}

// Go records VCS metadata for ordinary local builds when the checkout is available.
func buildCommit(info *debug.BuildInfo) string {
	if commit != "unknown" || info == nil {
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
