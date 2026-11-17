# Publish the first release

Goal: publish a versioned GitHub release for the standalone layerctl repository, with tag-triggered builds and installation artifacts. [Owners](../../OWNERS.md).

- /root: completed. GitHub and the SSH remote have no releases or tags; `main` is `3034ffcf7284d49112c703707de65b7141acea97` and its CI passed. The existing workflow builds four Darwin/Linux architectures with GoReleaser on version tags.
- Prepare the first semantic version, `v0.1.0`, and dated changelog notes. Verify the release notes, Go tests, four target builds, and workflow configuration before publishing.
- Push the release commit and tag through SSH, monitor GitHub Actions, and verify the release notes, version metadata, archives, and checksums. Preserve unrelated work and clean up temporary build artifacts.
- The Nix package still reports `0.0.0-unstable-20261114`; align its version with `0.1.0` and document installation from the tag and release assets.

- 2026-11-17: Go tests, release-note tests, GoReleaser configuration and tagged dry run passed. All four archives passed SHA-256 checks; the macOS ARM64 archive reports version and the test commit. Nix flake evaluation passed for all four systems and the native Nix package built as `0.1.0`. An offline Nix build failed because a required source was unavailable; the subsequent cached online build passed. Next: checkpoint, push `main` and `v0.1.0`, and verify the live release.

- 2026-11-17: `v0.1.0` released with four archives and checksums; both release and CI workflows passed. A real `go install ...@v0.1.0` reported `0.0.0-dev (unknown)`. Fix version reporting for module installs without changing embedded release versions; test it, bump Nix and documentation to `v0.1.1`, and publish and verify the patch release. Preserve the published tag.

- 2026-11-17: Patch validation passed: the full Go suite and release-note checks passed, an isolated Go module proxy produced `layerctl 0.1.1 (unknown)` from `go install @v0.1.1`, and the native Nix package built and reported `0.1.1`. Next: checkpoint and push the patch tag, then verify GitHub assets and public Go installation.

- 2026-11-17: Published `v0.1.1` at commit `f14eb051eb484a5583cc3e1374450a271cc01886`. [Release workflow](https://github.com/NonLogicalDev/util.layerctl/actions/runs/44146956202) and main CI passed. All four published archives passed checksum verification; the macOS ARM64 executable reported the exact version and commit. Public `go install @v0.1.1` reported version `0.1.1`; Nix built the package. `v0.1.0` remains published as the initial release.
