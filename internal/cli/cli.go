// Modified from abg-OAI/codex: module import path changed for this extraction.
// Package cli parses layerctl commands and renders their results.
// Domain packages own repository, projection, and layer operations.
package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
)

// Run parses and executes one invocation. getwd and version are called only
// when the selected command needs them.
func Run(
	ctx context.Context,
	args []string,
	logger *log.Logger,
	getwd func() (string, error),
	version func() string,
	options ...kong.Option,
) error {
	parser, err := kong.New(new(rootCmd), append([]kong.Option{
		kong.Name("layerctl"),
		kong.Description("Compose Git repositories and changes as ordered layers."),
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.Bind(&repositoryLoader{Getwd: getwd, Log: logger}, versionProvider(version)),
	}, options...)...)
	if err != nil {
		return fmt.Errorf("construct CLI: %w", err)
	}
	parsed, err := parser.Parse(args)
	if err != nil {
		return err
	}
	return parsed.Run()
}

type rootCmd struct {
	VersionFlag versionFlag `name:"version" help:"Print the version and source commit. Use without a command."`

	Projection projectionCmd `cmd:"" help:"Manage projection workspaces."`
	Layer      layerCmd      `cmd:"" help:"Capture changes as layers."`
	Graph      graphCmd      `cmd:"" help:"Render the layer dependency graph."`
	Check      checkCmd      `cmd:"" help:"Check that the definition composes successfully."`
	Version    versionCmd    `cmd:"" help:"Print the version and source commit."`
}
