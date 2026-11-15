# Repository Guidelines

## Project layout

Layerctl composes Git sources and ordered operations into projection workspaces. The CLI entrypoint is `cmd/layerctl/`; `internal/` contains definition parsing, Git operations, projection, layer capture, and graph rendering. Tests live alongside the packages they exercise. `skills/layerctl/` contains the user-facing authoring guide. The project retains attribution to the original OpenAI extraction in `README.md` and `LICENSE`.

## Development and validation

Go 1.26 and Git are required. Run `go test ./...` and `go build -o /dev/null ./cmd/layerctl` from the repository root. Run `gofmt` on changed Go files. Tests create disposable Git repositories; preserve existing projections and worktrees when debugging. Use `nix flake check` to build and test the Nix package, and `goreleaser check` to validate the release configuration. Run `./scripts/test-release-notes.sh` after changing release-note extraction. `goreleaser release --snapshot --clean` builds release artifacts without publishing them.

## Packaging and releases

`nix/package.nix` holds the reusable `buildGoModule` definition; `flake.nix` exposes it for Linux and macOS on x86-64 and ARM64. Keep the vendored Go dependency hash current. Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) in `CHANGELOG.md`: record notable changes under `Unreleased`, then prepare a dated semantic-version section before tagging. The release workflow reads that section and runs GoReleaser on a `v`-prefixed version tag. Release builds embed the version and full commit in `layerctl --version`. Do not create or push a release tag without a request.

## Changes and review

Keep changes focused and preserve upstream attribution. Read applicable instructions under `docs/plans/` before changing project plans. Explain observable behavior, relevant validation, and limitations in commits or reviews. Run focused checks and the full Go suite before handing off code changes; do not claim a check passed unless it ran.
