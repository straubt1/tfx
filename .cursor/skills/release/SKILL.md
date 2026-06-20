---
name: release
description: >-
  Prepare and cut TFx semver releases (CHANGELOG, validation, tagging via Taskfile).
  Use when the user asks to prepare a release, update CHANGELOG for a new version,
  cut vX.Y.Z, run release:patch/minor/major, or avoid post-merge release failures.
disable-model-invocation: true
---

# TFx Release

Two-phase workflow: **prepare on a branch** (PR) → **tag on `main` after merge**. Never tag before the CHANGELOG PR is merged.

## Phase 1 — Prepare (branch + PR)

### 1. Choose the version bump

| Bump | Task | Example |
|------|------|---------|
| Patch — bugfixes only | `task release:patch` | v0.4.0 → v0.4.1 |
| Minor — new features | `task release:minor` | v0.3.3 → v0.4.0 |
| Major — breaking changes | `task release:major` | v0.4.0 → v1.0.0 |

Confirm the target version with the user when ambiguous.

### 2. Create a release branch

```bash
git checkout main && git pull
git checkout -b release/vX.Y.Z
```

### 3. Update CHANGELOG.md

This repo uses [Keep a Changelog](https://keepachangelog.com/) with versioned unreleased sections:

```markdown
## [vX.Y.Z] - Unreleased
```

**If an unreleased section already exists** (e.g. `[v0.3.4] - Unreleased` but releasing v0.4.0): rename the heading to the target version and review bullets.

**If no unreleased section exists**, create one after the intro block. Gather changes since the latest tag:

```bash
git log $(git describe --tags --abbrev=0 --match "v[0-9]*.[0-9]*.[0-9]*")..HEAD --oneline
```

Group bullets under **Added**, **Changed**, **Fixed**, **Removed**. Include PR numbers where available.

**Fix compare links** at the bottom of `CHANGELOG.md`:

```markdown
[Unreleased]: https://github.com/straubt1/tfx/compare/vX.Y.Z...HEAD
[vX.Y.Z]: https://github.com/straubt1/tfx/compare/vPREVIOUS...vX.Y.Z
```

Replace `vPREVIOUS` with the tag before this release.

Leave the date as `Unreleased` in the prep PR.

### 4. Validate before opening the PR

Run all checks — do not skip:

```bash
# CHANGELOG must contain the exact version string the release task will look for
grep "\[vX.Y.Z\]" CHANGELOG.md

go build ./...
go test ./...

# Simulates goreleaser build (no tag required); cask written to dist/
task release:dry-run
```

Goreleaser v2 may emit 2 informational `dockers_v2` warnings; these are harmless.

**Do not edit `version/version.go`** — goreleaser injects the version from the git tag at build time.

### 5. Open the PR

Stage only release-prep files (typically `CHANGELOG.md`; include the skill if updated):

```bash
git add CHANGELOG.md
git push -u origin HEAD
gh pr create --title "Prepare release vX.Y.Z" --body "$(cat <<'EOF'
## Summary
- Finalize CHANGELOG for vX.Y.Z

## Test plan
- [ ] `grep "\[vX.Y.Z\]" CHANGELOG.md` passes
- [ ] `go build ./...` and `go test ./...` pass
- [ ] `task release:dry-run` succeeds

## After merge
On `main`, run the matching release task: `task release:<patch|minor|major>`.
EOF
)"
```

Do not commit or push unless the user asks.

## Phase 2 — Cut the release (after PR merge)

**Use the Taskfile release tasks only.** Do not tag, push tags, or cut the release manually.

On up-to-date `main`:

```bash
git checkout main && git pull
task release:minor   # pick the task that matches the bump chosen in Phase 1
```

| Bump | Task |
|------|------|
| Patch | `task release:patch` |
| Minor | `task release:minor` |
| Major | `task release:major` |

Each task delegates to `_release-bump`, which:

- Computes the next semver from the latest `v*.*.*` git tag
- Warns if `CHANGELOG.md` lacks a `[vNEXT_VERSION]` entry
- Prompts for confirmation, then tags and pushes

After the task completes, monitor [GitHub Actions](https://github.com/straubt1/tfx/actions) — the release workflow runs on tag push and publishes binaries, Docker (GHCR), Linux packages, and the Homebrew tap.

Requires `HOMEBREW_TAP_TOKEN` in repo secrets (classic PAT with `repo` scope on `straubt1/homebrew-tap`).

To test the pipeline locally before tagging: `task release:dry-run`.

## Common post-merge failures

| Symptom | Cause | Fix |
|---------|-------|-----|
| `_release-bump` warns about missing CHANGELOG | Heading uses wrong version (e.g. `v0.3.4` but tagging `v0.4.0`) | Rename section to exact target: `## [v0.4.0]` |
| Goreleaser fails in CI | Build/test issue on `main` | Fix on `main`, delete bad tag if created, re-tag |
| Homebrew tap not updated | Missing/expired `HOMEBREW_TAP_TOKEN` | Check repo secrets; re-run release workflow |
| Wrong version published | Used wrong bump task | Delete remote tag (with user approval), fix CHANGELOG, re-cut |

## Quick reference

```bash
git describe --tags --abbrev=0 --match "v[0-9]*.[0-9]*.[0-9]*"   # current tag
task release:dry-run                                               # local pipeline test
task release:patch | release:minor | release:major                   # cut release
```

See also: [DEVELOPMENT.md § Releasing](../../DEVELOPMENT.md), [CLAUDE.md § Release Process](../../CLAUDE.md).
