package cli

import (
	"errors"
	"fmt"

	"github.com/alecthomas/kong"
)

type versionCmd struct{}

// Run prints the version and source commit.
func (*versionCmd) Run(kctx *kong.Context, version versionProvider) error {
	_, err := fmt.Fprintln(kctx.Stdout, "layerctl", version())
	return err
}

type versionProvider func() string

type versionFlag bool

// BeforeReset handles --version before Kong requires a subcommand.
func (versionFlag) BeforeReset(kctx *kong.Context, path *kong.Path, version versionProvider) error {
	if !kctx.Value(path).Bool() {
		return nil
	}
	if kctx.Selected() != nil {
		return errors.New("--version does not accept a command")
	}
	if err := new(versionCmd).Run(kctx, version); err != nil {
		return err
	}
	kctx.Kong.Exit(0)
	return nil
}
