// Modified from abg-OAI/codex: module import path changed for this extraction.
package cli

import (
	"context"
	"errors"

	"github.com/alecthomas/kong"

	"github.com/NonLogicalDev/util.layerctl/internal/layer"
	"github.com/NonLogicalDev/util.layerctl/internal/projection"
)

type layerCmd struct {
	Add     layerAddCmd     `cmd:"" help:"Capture a new layer from a projection."`
	Refresh layerRefreshCmd `cmd:"" help:"Refresh an existing layer from a projection."`
}

// AfterApply binds a shared composer for capture and refresh after Kong
// validates the invocation. The providers load the repository on demand.
func (*layerCmd) AfterApply(kctx *kong.Context, repo *repositoryLoader) error {
	return errors.Join(
		kctx.BindSingletonProvider(repo.OpenProjection),
		kctx.BindSingletonProvider(func(composer *projection.CompositionService) *layer.Service {
			return &layer.Service{
				Definition: composer.Definition,
				Git:        composer.Git,
				Projection: composer,
			}
		}),
	)
}

type layerAddCmd struct {
	ID     string         `arg:"" help:"New layer ID (NNNN-LAYER-SLUG)."`
	Source layerSourceArg `embed:""`
}

// Run captures the projection's changes as a new layer.
func (c *layerAddCmd) Run(ctx context.Context, service *layer.Service) error {
	return service.Add(ctx, layer.AddRequest{ID: c.ID, Projection: c.Source.From})
}

func (*layerAddCmd) Help() string {
	return "Create a projection through the preceding enabled layer " +
		"and commit the changes to capture in its workspace."
}

type layerRefreshCmd struct {
	ID     string         `arg:"" help:"Layer ID."`
	Source layerSourceArg `embed:""`
}

// Run updates an existing layer from the projection's changes.
func (c *layerRefreshCmd) Run(ctx context.Context, service *layer.Service, composer *projection.CompositionService) error {
	return service.RefreshComposition(ctx, c.ID, c.Source.From, composer)
}

func (*layerRefreshCmd) Help() string {
	return "Create a projection through this layer " +
		"and commit the desired changes in its workspace. " +
		"Refresh preserves the existing operations and appends a patch. " +
		"Verification replays the operations, including scripts."
}

type layerSourceArg struct {
	From string `required:"" placeholder:"PROJECTION" help:"Capture from PROJECTION."`
}

// Validate requires a nonempty projection name for --from.
func (arg *layerSourceArg) Validate() error {
	if arg.From == "" {
		return errors.New("--from must name a projection")
	}
	return nil
}
