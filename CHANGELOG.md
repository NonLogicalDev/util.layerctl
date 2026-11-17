# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-11-17

### Added

- Import repositories containing Git submodules with their pins and original
  `.gitmodules` file. Submodules remain uninitialized.
- Show usage and available flags with `layerctl --help` and `--help` on subcommands.

### Changed

- Accept command flags before or after positional arguments.

## [0.1.1] - 2026-11-17

### Fixed

- Report the module version for binaries installed with `go install` at a version tag.

## [0.1.0] - 2026-11-17

### Added

- Compose one or more Git repositories with ordered clone, patch, overlay, and script operations into Git projection workspaces.
- Preserve imported repository history as merge parents and prepare source trees with pre-checkout hooks.
- Define layer dependencies and commit metadata in TOML, render DOT or Mermaid graphs, and capture or refresh changes from a projection.
- Build and run Layerctl with Nix on supported Linux and macOS systems.
- Report the executable's embedded version and commit with `layerctl --version`.
- Publish versioned archives and checksums with GoReleaser when a release tag is pushed.

### Changed

- Require manifest-based layers; the original `upstream.json` format and its migration commands are not supported.

[Unreleased]: https://github.com/NonLogicalDev/util.layerctl/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/NonLogicalDev/util.layerctl/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/NonLogicalDev/util.layerctl/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/NonLogicalDev/util.layerctl/releases/tag/v0.1.0
