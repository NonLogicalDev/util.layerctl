// Modified from abg-OAI/codex: module import path changed for this extraction.
package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/NonLogicalDev/util.layerctl/internal/definition"
	"github.com/NonLogicalDev/util.layerctl/internal/gitrepo"
	"github.com/NonLogicalDev/util.layerctl/internal/projection"
)

type repositoryLoader struct {
	Getwd func() (string, error) // required
	Log   *log.Logger            // required
}

// Load discovers the Git repository containing the working directory and loads
// its layer definition.
func (r *repositoryLoader) Load(ctx context.Context) (*definition.Repository, *gitrepo.Repository, error) {
	dir, err := r.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("get working directory: %w", err)
	}
	git, err := gitrepo.Discover(ctx, dir)
	if err != nil {
		return nil, nil, fmt.Errorf("discover Git repository: %w", err)
	}
	canonical, err := definition.Load(git.Root)
	if err != nil {
		return nil, nil, fmt.Errorf("load canonical definition: %w", err)
	}
	return canonical, git, nil
}

// OpenProjection loads the definition and constructs a projection service
// that uses the invocation logger for progress.
func (r *repositoryLoader) OpenProjection(ctx context.Context) (*projection.CompositionService, error) {
	canonical, git, err := r.Load(ctx)
	if err != nil {
		return nil, err
	}
	return &projection.CompositionService{
		Definition: canonical,
		Git:        git,
		Log:        r.Log,
	}, nil
}
