package cli

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
)

func TestRun_usageErrors(t *testing.T) {
	tests := []struct {
		name     string
		giveArgs []string
		want     string
	}{
		{"MissingCommand", nil, "expected one of"},
		{"MissingProjectionSubcommand", []string{"projection"}, "expected one of"},
		{"ExtraProjectionName", []string{"projection", "path", "demo", "extra"}, "unexpected argument extra"},
		{"MissingLayerSource", []string{"layer", "refresh", "0001-base"}, "--from must name a projection"},
		{"EmptyLayerSource", []string{"layer", "add", "0002-feature", "--from="}, "--from must name a projection"},
		{"InvalidGraphFormat", []string{"graph", "--format", "html"}, "one of"},
		{"VersionWithCommand", []string{"--version", "graph"}, "--version does not accept a command"},
		{"UnexpectedGraphArgument", []string{"graph", "extra"}, "unexpected argument extra"},
		{"UnknownGraphFlag", []string{"graph", "--unknown"}, "unknown flag --unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			err := Run(t.Context(), tt.giveArgs, log.New(&stderr, "", 0),
				func() (string, error) {
					t.Fatal("invalid arguments must not query the working directory")
					return "", nil
				},
				func() string {
					t.Fatal("invalid arguments must not query the version")
					return ""
				},
				kong.Exit(func(code int) { t.Fatalf("invalid arguments requested exit %d", code) }),
			)
			if _, ok := errors.AsType[*kong.ParseError](err); !ok || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Run(%v) error = %v, want a parse error containing %q", tt.giveArgs, err, tt.want)
			}
			if stderr.Len() != 0 {
				t.Errorf("Run(%v) stderr = %q, want empty output before reporting the error", tt.giveArgs, &stderr)
			}
		})
	}
}

func TestRun_queriesProcessStateOnDemand(t *testing.T) {
	getwdErr := errors.New("working directory unavailable")
	tests := []struct {
		name             string
		giveArgs         []string
		wantErr          error
		wantExit         bool
		wantGetwdCalls   int
		wantVersionCalls int
	}{
		{name: "RootHelp", giveArgs: []string{"--help"}, wantExit: true},
		{name: "ProjectionPathHelp", giveArgs: []string{"projection", "path", "demo", "--help"}, wantExit: true},
		{name: "VersionHelp", giveArgs: []string{"version", "--help"}, wantExit: true},
		{name: "Graph", giveArgs: []string{"graph"}, wantErr: getwdErr, wantGetwdCalls: 1},
		{name: "LayerRefresh", giveArgs: []string{"layer", "refresh", "0001-base", "--from", "demo"}, wantErr: getwdErr, wantGetwdCalls: 1},
		{name: "VersionFlag", giveArgs: []string{"--version"}, wantExit: true, wantVersionCalls: 1},
		{name: "VersionCommand", giveArgs: []string{"version"}, wantVersionCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var getwdCalls, versionCalls int
			exited, err := runUntilExit(t, func(exit kong.Option) error {
				return Run(t.Context(), tt.giveArgs, log.New(io.Discard, "", 0),
					func() (string, error) {
						getwdCalls++
						return "", getwdErr
					},
					func() string {
						versionCalls++
						return "test"
					},
					exit,
				)
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Run(%v) error = %v, want %v", tt.giveArgs, err, tt.wantErr)
			}
			if exited != tt.wantExit {
				t.Errorf("Run(%v) called Exit = %v, want %v", tt.giveArgs, exited, tt.wantExit)
			}
			if getwdCalls != tt.wantGetwdCalls || versionCalls != tt.wantVersionCalls {
				t.Errorf("Getwd called %d times and Version called %d times, want %d and %d",
					getwdCalls, versionCalls, tt.wantGetwdCalls, tt.wantVersionCalls)
			}
		})
	}
}

func TestProjectionCreateFlagOrder(t *testing.T) {
	for _, args := range [][]string{
		{"projection", "create", "demo", "--worktree", "work", "--through=0001-base"},
		{"projection", "create", "--worktree", "work", "--through=0001-base", "demo"},
	} {
		var tree rootCmd
		parser, err := kong.New(&tree)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := parser.Parse(args); err != nil {
			t.Fatalf("Parse(%v): %v", args, err)
		}
		want := projectionCreateCmd{Name: "demo", Worktree: "work", Through: "0001-base"}
		if got := tree.Projection.Create; got != want {
			t.Errorf("Parse(%v) create command = %+v, want %+v", args, got, want)
		}
	}
}

func runUntilExit(t *testing.T, run func(kong.Option) error) (exited bool, err error) {
	t.Helper()
	type exitStatus int
	defer func() {
		if value := recover(); value != nil {
			code, ok := value.(exitStatus)
			if !ok {
				panic(value)
			}
			if code != 0 {
				t.Errorf("Run requested exit %d, want zero", code)
			}
			exited = true
		}
	}()
	return false, run(kong.Exit(func(code int) { panic(exitStatus(code)) }))
}
