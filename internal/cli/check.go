// Modified from abg-OAI/codex: module import path changed for this extraction.
package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NonLogicalDev/util.layerctl/internal/projection"
)

type checkCmd struct{}

func (*checkCmd) Help() string {
	return "The check runs the definition's scripts. " +
		"If composition fails, " +
		"the temporary projection can remain available for inspection."
}

// Run checks the definition by building a temporary projection.
func (*checkCmd) Run(ctx context.Context, repo *repositoryLoader) error {
	service, err := repo.OpenProjection(ctx)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "layerctl-check-")
	if err != nil {
		return err
	}
	name := strings.ToLower(filepath.Base(f.Name()))
	if err := f.Close(); err != nil {
		return err
	}
	_ = os.Remove(f.Name())
	_, err = service.Create(ctx, projection.CreateRequest{Name: name})
	if err != nil {
		return fmt.Errorf("check composition: %w", err)
	}
	if err := service.Delete(ctx, name); err != nil {
		return fmt.Errorf("remove check projection: %w", err)
	}
	return nil
}
