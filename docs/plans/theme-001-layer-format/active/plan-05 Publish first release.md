# Publish the first release

Goal: publish a versioned GitHub release for the standalone layerctl repository, with tag-triggered builds and installation artifacts. [Owners](../../OWNERS.md).

- /root: active. GitHub and the SSH remote have no releases or tags; `main` is `3034ffcf7284d49112c703707de65b7141acea97` and its CI passed. The existing workflow builds four Darwin/Linux architectures with GoReleaser on version tags.
- Prepare the first semantic version, `v0.1.0`, and dated changelog notes. Verify the release notes, Go tests, four target builds, and workflow configuration before publishing.
- Push the release commit and tag through SSH, monitor GitHub Actions, and verify the release notes, version metadata, archives, and checksums. Preserve unrelated work and clean up temporary build artifacts.
- The Nix package still reports `0.0.0-unstable-20261114`; align its version with `0.1.0` and document installation from the tag and release assets.

- 2026-11-17: Go tests, release-note tests, GoReleaser configuration and tagged dry run passed. All four archives passed SHA-256 checks; the macOS ARM64 archive reports version and the test commit. Nix flake evaluation passed for all four systems and the native Nix package built as `0.1.0`. An offline Nix build failed because a required source was unavailable; the subsequent cached online build passed. Next: checkpoint, push `main` and `v0.1.0`, and verify the live release.
