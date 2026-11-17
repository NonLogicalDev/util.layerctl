// Modified from abg-OAI/codex: module import path changed for this extraction.
package cli

import (
	"context"

	"github.com/alecthomas/kong"

	"github.com/NonLogicalDev/util.layerctl/internal/graph"
)

type graphCmd struct {
	Format string `default:"dot" enum:"dot,mermaid" help:"Graph output format (${enum})."`
}

// Run renders the layer dependency graph.
func (c *graphCmd) Run(ctx context.Context, repo *repositoryLoader, kctx *kong.Context) error {
	canonical, _, err := repo.Load(ctx)
	if err != nil {
		return err
	}
	return graph.Write(kctx.Stdout, canonical.Layers, c.Format)
}
