package main

import (
	"os/exec"
	"path/filepath"
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
