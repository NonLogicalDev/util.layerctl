// Modified from abg-OAI/codex: module import path changed for this extraction.
package cli

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"

	"github.com/NonLogicalDev/util.layerctl/internal/projection"
)

type projectionCmd struct {
	Create   projectionCreateCmd   `cmd:"" help:"Create a projection workspace."`
	List     projectionListCmd     `cmd:"" help:"List projections."`
	Path     projectionPathCmd     `cmd:"" help:"Print the path to a projection workspace."`
	Continue projectionContinueCmd `cmd:"" help:"Resume a projection after resolving a conflict."`
	Unlock   projectionUnlockCmd   `cmd:"" help:"Clear a stale projection lock."`
	Delete   projectionDeleteCmd   `cmd:"" help:"Remove a projection and its workspace."`
}

type projectionCreateCmd struct {
	Name     string `arg:"" help:"Projection name."`
	Worktree string `placeholder:"PATH" help:"Attach the projection at PATH."`
	Through  string `placeholder:"LAYER" help:"Stop after LAYER."`
}

func (*projectionCreateCmd) Help() string {
	return "The worktree path must not already exist. " +
		"If you omit --worktree, " +
		"Layerctl creates the workspace in its local Git state directory."
}

// Run creates a projection and prints its workspace path.
func (c *projectionCreateCmd) Run(ctx context.Context, repo *repositoryLoader, kctx *kong.Context) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	created, err := service.Create(ctx, projection.CreateRequest{
		Name:         c.Name,
		WorktreePath: c.Worktree,
		Through:      c.Through,
	})
	if err != nil {
		return err
	}
	fmt.Fprintln(kctx.Stdout, created.WorktreePath)
	return nil
}

type projectionListCmd struct{}

// Run prints one row per projection.
func (*projectionListCmd) Run(ctx context.Context, repo *repositoryLoader, kctx *kong.Context) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	projections, err := service.List(ctx)
	if err != nil {
		return err
	}
	for _, item := range projections {
		fmt.Fprintf(kctx.Stdout, "%s\t%s\t%s\t%s\n", item.Name, item.Base, item.Head, item.WorktreePath)
	}
	return nil
}

type projectionPathCmd struct {
	Name string `arg:"" help:"Projection name."`
}

// Run prints a projection's workspace path.
func (c *projectionPathCmd) Run(ctx context.Context, repo *repositoryLoader, kctx *kong.Context) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	path, err := service.Path(ctx, c.Name)
	if err != nil {
		return err
	}
	fmt.Fprintln(kctx.Stdout, path)
	return nil
}

type projectionContinueCmd struct {
	Name string `arg:"" help:"Projection name."`
}

func (*projectionContinueCmd) Help() string {
	return "Resolve the conflict in the reported workspace before continuing. " +
		"If the process stopped during an operation, delete and recreate the projection."
}

// Run resumes a projection and reports its workspace path when available.
func (c *projectionContinueCmd) Run(ctx context.Context, repo *repositoryLoader, kctx *kong.Context) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	result, err := service.Continue(ctx, c.Name)
	if result.WorktreePath != "" {
		fmt.Fprintln(kctx.Stdout, result.WorktreePath)
	}
	return err
}

type projectionUnlockCmd struct {
	Name string `arg:"" help:"Projection name."`
}

func (*projectionUnlockCmd) Help() string {
	return "Verify that no Layerctl command is operating on the projection " +
		"before clearing its lock."
}

// Run clears a stale projection lock.
func (c *projectionUnlockCmd) Run(ctx context.Context, repo *repositoryLoader) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	return service.Unlock(ctx, c.Name)
}

type projectionDeleteCmd struct {
	Name string `arg:"" help:"Projection name."`
}

func (*projectionDeleteCmd) Help() string {
	return "Deleting a projection also removes any uncaptured changes in its workspace."
}

// Run removes a projection and its workspace.
func (c *projectionDeleteCmd) Run(ctx context.Context, repo *repositoryLoader) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}
	return service.Delete(ctx, c.Name)
}
