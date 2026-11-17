package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	binary := buildLayerctl(t)
	dir := t.TempDir()
	tests := []struct {
		name     string
		giveArgs []string
		want     string
	}{
		{"Root", []string{"--help"}, "projection create <name>"},
		{"Short", []string{"-h"}, "Usage: layerctl"},
		{"Projection", []string{"projection", "--help"}, "projection unlock <name>"},
		{"ProjectionCreate", []string{"projection", "create", "--help"}, "--worktree=PATH"},
		{"ProjectionCreateWithName", []string{"projection", "create", "demo", "--help"}, "--through=LAYER"},
		{"ProjectionList", []string{"projection", "list", "--help"}, "Usage: layerctl projection list"},
		{"ProjectionPath", []string{"projection", "path", "--help"}, "Usage: layerctl projection path <name>"},
		{"ProjectionContinue", []string{"projection", "continue", "--help"}, "Usage: layerctl projection continue <name>"},
		{"ProjectionUnlock", []string{"projection", "unlock", "--help"}, "Usage: layerctl projection unlock <name>"},
		{"ProjectionDelete", []string{"projection", "delete", "--help"}, "Usage: layerctl projection delete <name>"},
		{"Layer", []string{"layer", "--help"}, "layer refresh --from=PROJECTION <id>"},
		{"LayerAdd", []string{"layer", "add", "--help"}, "--from=PROJECTION"},
		{"LayerRefresh", []string{"layer", "refresh", "--help"}, "--from=PROJECTION"},
		{"Graph", []string{"graph", "--help"}, "--format="},
		{"Check", []string{"check", "--help"}, "Usage: layerctl check"},
		{"Version", []string{"version", "--help"}, "Usage: layerctl version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.CommandContext(t.Context(), binary, tt.giveArgs...)
			cmd.Dir = dir
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			stdout, err := cmd.Output()
			if err != nil {
				t.Fatalf("layerctl %v: %v\n%s", tt.giveArgs, err, &stderr)
			}
			if !strings.Contains(string(stdout), tt.want) {
				t.Errorf("stdout = %q, want it to contain %q", stdout, tt.want)
			}
			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty output", &stderr)
			}
		})
	}
}

// The version command must work outside a definition repository and report
// the values embedded by release builds.
func TestReleaseVersion(t *testing.T) {
	binary := buildLayerctl(t, "-ldflags", "-X main.version=1.2.3 -X main.commit=abc123")
	for _, arg := range []string{"--version", "version"} {
		cmd := exec.CommandContext(t.Context(), binary, arg)
		cmd.Dir = t.TempDir()
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s: %v\n%s", arg, err, output)
		}
		if got, want := string(output), "layerctl 1.2.3 (abc123)\n"; got != want {
			t.Errorf("%s: got %q, want %q", arg, got, want)
		}
	}
}

func TestBuildVersion(t *testing.T) {
	for _, test := range []struct {
		name, embedded, module, want string
	}{
		{"tagged module", "0.0.0-dev", "v0.1.1", "0.1.1"},
		{"pseudo-version", "0.0.0-dev", "v0.0.0-20261117000000-abcdef123456", "0.0.0-20261117000000-abcdef123456"},
		{"local checkout", "0.0.0-dev", "(devel)", "0.0.0-dev"},
		{"embedded takes precedence", "1.2.3", "v0.1.1", "1.2.3"},
	} {
		t.Run(test.name, func(t *testing.T) {
			info := &debug.BuildInfo{Main: debug.Module{Version: test.module}}
			if got := buildVersion(test.embedded, info); got != test.want {
				t.Fatalf("buildVersion() = %q, want %q", got, test.want)
			}
		})
	}
	if got := buildVersion("0.0.0-dev", nil); got != "0.0.0-dev" {
		t.Fatalf("buildVersion() without build info = %q", got)
	}
}

func buildLayerctl(t *testing.T, buildFlags ...string) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "layerctl")
	args := append([]string{"build", "-o", binary}, buildFlags...)
	build := exec.CommandContext(t.Context(), "go", append(args, ".")...)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build layerctl: %v\n%s", err, output)
	}
	return binary
}
