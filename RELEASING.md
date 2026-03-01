# Releasing TFx

This document describes how to cut a new release of TFx.

## Prerequisites

Before your first release, ensure the following GitHub Actions secret is set in this repository:

| Secret | Description |
|--------|-------------|
| `HOMEBREW_TAP_TOKEN` | Classic GitHub PAT with `repo` scope on [`straubt1/homebrew-tap`](https://github.com/straubt1/homebrew-tap). Required for goreleaser to commit the updated cask. |

## What Happens on Release

Pushing a `v*` tag triggers `.github/workflows/release.yml`, which uses goreleaser to:

1. Build binaries for macOS (amd64/arm64), Linux (amd64/arm64), and Windows (amd64/arm64)
2. Create archives and checksums
3. Publish a GitHub Release with all artifacts
4. Build and push a multi-platform Docker image to GHCR (`ghcr.io/straubt1/tfx`)
5. Package `.apk`, `.deb`, and `.rpm` Linux packages
6. Commit an updated cask to [`straubt1/homebrew-tap`](https://github.com/straubt1/homebrew-tap)

## Release Steps

### 1. Update CHANGELOG.md

Edit `CHANGELOG.md` and fill in the entry for the new version. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/):

```markdown
## [v0.2.0] - 2026-03-01

**Added**
* ...

**Changed**
* ...

**Fixed**
* ...
```

The release tasks will warn you if a CHANGELOG entry is missing for the target version.

### 2. Run the release task

Pick the appropriate bump type and run one of:

```bash
task release:patch   # x.y.Z+1  — bug fixes, no new features
task release:minor   # x.Y+1.0  — new features, backward-compatible
task release:major   # X+1.0.0  — breaking changes
```

Each task will:
- Detect the current version from the latest git tag
- Compute the next version
- Show a preview of what will happen
- Warn if `CHANGELOG.md` is missing an entry
- Prompt for confirmation before proceeding
- Commit `CHANGELOG.md`, create the tag, and push

### 3. Monitor the release

GitHub Actions will pick up the new tag automatically:
<https://github.com/straubt1/tfx/actions>

Once complete, verify:
- [ ] GitHub Release page has all expected artifacts
- [ ] `docker pull ghcr.io/straubt1/tfx:<tag>` works
- [ ] Homebrew tap has a new commit: <https://github.com/straubt1/homebrew-tap>
- [ ] `brew upgrade tfx` installs the new version

## Testing the Release Pipeline Locally

To validate the full build without publishing anything:

```bash
task release-dry-run
```

This runs `goreleaser release --snapshot --clean --skip=announce,validate`. In snapshot mode:
- No tag is required
- Nothing is published or pushed
- The Homebrew cask is generated and written to `dist/` instead of pushed to the tap (controlled by `skip_upload: "{{ .IsSnapshot }}"` in `.goreleaser.yml`)

After the run, inspect the generated cask:

```bash
cat dist/tfx.rb
```

> **Note:** goreleaser v2 emits 2 informational `dockers_v2` deprecation warnings during `goreleaser check` and dry runs. These are a known goreleaser quirk and can be ignored — the config is correct.

## Version Numbering

TFx follows [Semantic Versioning](https://semver.org/):

| Change | Example | When to use |
|--------|---------|-------------|
| Patch | `0.1.5` → `0.1.6` | Bug fixes, docs, internal refactors |
| Minor | `0.1.5` → `0.2.0` | New commands or flags, backward-compatible |
| Major | `0.1.5` → `1.0.0` | Breaking changes to CLI interface or config |

`version/version.go` defaults to `"dev"` and is never manually edited. Goreleaser injects the real version from the git tag at build time via ldflags.
