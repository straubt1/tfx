---
name: go-upgrade
description: Upgrade macOS Go and all module dependencies in TFx via task go:upgrade, then verify build and tests. Use when the user asks to upgrade Go, update dependencies, bump go.mod, or refresh packages.
disable-model-invocation: true
---

# Go Upgrade

Upgrade the local Go toolchain and all module dependencies, then verify the repo still builds and tests pass.

## Prerequisites

- macOS with Homebrew
- Run from repo root
- If tools may be missing: `task development:doctor` (install with `brew install go goreleaser go-task`)

## Upgrade

Run the Taskfile task (see [DEVELOPMENT.md](../../DEVELOPMENT.md) § Upgrading Dependencies):

```bash
task go:upgrade
```

This upgrades `go` and `goreleaser` via Homebrew, syncs `go.mod` to the installed Go version, runs `go get -u ./...`, `go mod tidy`, and updates the README Go version badge.

**Files touched:** `go.mod`, `go.sum`, `README.md` (badge only), `CHANGELOG.md`. CI reads Go version from `go.mod` via `go-version-file`.

If `brew upgrade go` reports already latest, the task still refreshes module dependencies.

## Verify

Match CI ([`.github/workflows/build.yml`](../../.github/workflows/build.yml)):

```bash
go build ./...
go test ./...
```

Some `client/` tests call the TFE API — run with network access if the sandbox blocks them.

On failure: diagnose and fix compile or test issues from breaking dependency bumps. Do not downgrade Go without asking the user.

Optional smoke check:

```bash
task go:build && ./tfx version
```

Integration tests are out of scope unless requested (`secrets/.env-int` required):

```bash
task test:integration-data
task test:integration-cmd
```

## Changelog

Update [`CHANGELOG.md`](../../CHANGELOG.md) after a successful upgrade. This repo uses [Keep a Changelog](https://keepachangelog.com/) with unreleased sections titled `## [vX.Y.Z] - Unreleased` (not a bare `[Unreleased]` heading).

1. Read the top of `CHANGELOG.md` and look for a section matching `## [vX.Y.Z] - Unreleased`.
2. **If it exists:** add a bullet under **Changed** (create a `**Changed**` subsection if missing):

   ```
   * Upgraded Go to X.Y.Z and refreshed module dependencies
   ```

   Replace `X.Y.Z` with the version from `go version`. Do not duplicate the entry if an equivalent bullet is already present.

3. **If it does not exist:** insert a new unreleased section immediately after the Keep a Changelog intro (before the first dated release). Derive the version label from the latest git tag:

   ```bash
   git describe --tags --abbrev=0 --match "v[0-9]*.[0-9]*.[0-9]*"
   ```

   Bump the patch segment (e.g. `v0.3.3` → `v0.3.4`) and add:

   ```markdown
   ## [v0.3.4] - Unreleased

   **Changed**

   * Upgraded Go to X.Y.Z and refreshed module dependencies
   ```

## Report

1. Show `git diff --stat` for `go.mod`, `go.sum`, `README.md`, and `CHANGELOG.md`
2. Draft commit message (do not commit unless asked):

```
chore: upgrade Go to X.Y.Z and refresh module dependencies
```

Replace `X.Y.Z` with the version from `go version`.

## Default: stop before git

Do **not** create a branch, commit, or open a PR unless the user explicitly asks.

## Optional: branch, commit, and PR

When the user asks to commit and/or open a PR (e.g. "upgrade go and open a PR"):

1. Create branch: `git checkout -b chore/go-upgrade-<version>`
2. Stage only: `go.mod`, `go.sum`, `README.md`, `CHANGELOG.md`
3. Commit with HEREDOC message (use draft above)
4. Push: `git push -u origin HEAD`
5. Open PR: `gh pr create` with Summary and Test plan (note local `go build ./...` and `go test ./...` passed)

Do not push or commit unless explicitly requested.
