package main

import (
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"testing"
)

// The version command must work outside a definition repository and report
// the values embedded by release builds.
func TestReleaseVersion(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "layerctl")
	build := exec.Command("go", "build", "-ldflags", "-X main.version=1.2.3 -X main.commit=abc123", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build release binary: %v\n%s", err, output)
	}
	for _, arg := range []string{"--version", "version"} {
		cmd := exec.Command(binary, arg)
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
