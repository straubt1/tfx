# CLAUDE.md

This file provides context for AI assistants working in this codebase.

## Project Overview

**TFx** is a standalone CLI tool for HCP Terraform (Terraform Cloud) and Terraform Enterprise (TFE). It provides an API-driven interface for common operations that would otherwise require direct REST API calls.

## Development Setup

Required tools (macOS):
```bash
brew install go
brew install goreleaser
brew install mkdocs
```

## Git Workflow

**Do not commit changes unless explicitly asked.** After implementing a task, leave changes staged or unstaged for the user to review. Only run `git commit` when the user says to commit.

After completing a task, compose and display the intended commit message describing what changed, but do not execute `git commit`. This lets the user copy or adjust the message when they're ready to commit.

**Do not create new branches.** Make all changes directly on the currently active branch.

## Build & Test Commands

```bash
# Build
task go-build

# Cross-platform build (via goreleaser)
task go-build-all

# Run unit tests
task test

# Run integration tests (requires secrets/.env-int with TFE_HOSTNAME, TFE_TOKEN, TFE_ORGANIZATION)
task test-integration-data
task test-integration-cmd

# Run all tests
task test-all

# Upgrade dependencies
task go-upgrade

# Serve documentation site
task serve-docs
```

## Configuration & Environment Variables

**Required for integration tests and CLI use:**
- `TFE_TOKEN` - API authentication token
- `TFE_ORGANIZATION` - Target organization name
- `TFE_HOSTNAME` - TFE instance hostname (default: `app.terraform.io`)

**Optional:**
- `TFX_LOG` - Enable HTTP request/response logging for debugging
- `--config` flag or `~/.tfx.hcl` / `./.tfx.hcl` - Config file in HCL format

Integration tests load secrets from `secrets/.env-int` (not in repo).

## Project Structure

```
cmd/          # Cobra command implementations and orchestration
  flags/      # Per-command flag structs and parsers
  views/      # Output/rendering for each command
client/       # TFE API client wrapper (wraps go-tfe)
data/         # Data fetching layer — API calls and business logic
output/       # Output system (terminal tables, JSON, spinner, logger)
integration/  # Integration tests (require live TFE instance)
pkg/file/     # File utility package
version/      # Version string
main.go       # Entry point
Taskfile.yml  # Task runner configuration
```

## Architecture & Patterns

### Layered Architecture
1. **cmd/** — Command handlers, flag parsing, orchestration; calls data layer
2. **data/** — API calls and business logic using `TfxClient`
3. **client/** — Wraps `go-tfe` with context, HTTP logging, pagination helpers
4. **output/** — Rendering (terminal tables via `go-pretty`, JSON); never mixed with business logic
5. **flags/** — Flag struct definitions and `Parse*Flags(cmd)` helper functions

`cmd/project.go` is the reference implementation — follow its patterns for all new commands.

### Command Flow (per operation)
1. Parse flags → typed config struct (`cmd/flags/`)
2. Get client → `client.NewFromViper()`
3. Fetch data → function in `data/`
4. Render → view in `cmd/views/`

### Adding a New Command
To add e.g. `tfx project create`:
1. **`cmd/flags/project.go`** — Add `ProjectCreateFlags` struct and `ParseProjectCreateFlags(cmd)`
2. **`data/projects.go`** — Add `CreateProject(c *client.TfxClient, ...) (*tfe.Project, error)`
3. **`cmd/views/project_create.go`** — Add view to render the result
4. **`cmd/project.go`** — Wire it together in ~15 lines

### Before/After: Refactored Command Style

**Before (old monolithic style):**
```go
func projectList(cmd *cobra.Command) error {
    search, _ := cmd.Flags().GetString("search")
    client, _ := tfe.NewClient(config)
    pageNum := 1
    var all []*tfe.Project
    for {
        result, _ := client.Projects.List(ctx, org, &opts)
        all = append(all, result.Items...)
        if pageNum >= result.TotalPages { break }
        pageNum++
    }
    // inline rendering...
}
```

**After (refactored):**
```go
func projectList(cmdConfig *flags.ProjectListFlags) error {
    c, _ := client.NewFromViper()
    projects, _ := data.FetchProjects(c, c.OrganizationName, cmdConfig.Search)
    view := views.NewProjectListView(viper.GetBool("json"))
    return view.Render(c.OrganizationName, projects)
}
```

### Cobra Command Pattern
```go
var cmdName = &cobra.Command{
    Use:   "subcommand",
    Short: "Short description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // parse flags, get client, call data layer, render via views
    },
}
```

### Client Usage
```go
c, err := client.NewFromViper()  // reads Viper config/flags
```

### Error Handling
- Use `github.com/pkg/errors` — wrap errors with context via `errors.Wrap(err, "message")`
- View layer handles rendering errors with `RenderError()`

### Output Pattern
- Use `output.Get()` singleton for the shared output system
- View types expose `Render()` and `RenderError()` methods
- JSON output toggled via `-j / --json` flag

### Flag Pattern
- Define flag structs in `cmd/flags/`
- Parse with `Parse*Flags(cmd)` returning the struct
- Bind to Viper for config file / env var fallback

## Key Dependencies

- `github.com/hashicorp/go-tfe` — TFE/HCP Terraform API client
- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — Config and flag management
- `github.com/jedib0t/go-pretty` — Terminal table formatting
- `github.com/fatih/color` — Colored terminal output
- `github.com/pkg/errors` — Error wrapping

## Release Process

Releases are automated via goreleaser and GitHub Actions. The release workflow triggers on a new git tag and publishes binaries, a Docker image (GHCR), Linux packages (apk/deb/rpm), and a Homebrew cask update to `straubt1/homebrew-tap`.

### Prerequisites (one-time setup)

- `HOMEBREW_TAP_TOKEN` — classic GitHub PAT with `repo` scope on `straubt1/homebrew-tap`, stored as a GitHub Actions secret in this repo.

### Cutting a Release

Use one of these tasks — they auto-detect the current tag, compute the next semver, show a preview, check that CHANGELOG.md has an entry, then commit/tag/push:

```bash
task release:patch   # x.y.Z+1  — bugfixes
task release:minor   # x.Y+1.0  — new features
task release:major   # X+1.0.0  — breaking changes
```

Before confirming, update `CHANGELOG.md` with release notes for the new version.

### Testing the Release Pipeline Locally

```bash
task release-dry-run
```

Runs `goreleaser release --snapshot --clean --skip=announce,validate` — builds all artifacts without requiring a tag. The Homebrew cask is generated and written to `dist/` (not pushed to the tap) because `skip_upload: "{{ .IsSnapshot }}"` is set in `.goreleaser.yml`. Goreleaser v2 may emit 2 informational `dockers_v2` warnings; these are a known goreleaser quirk and can be ignored.

### Version Numbering

`version/version.go` defaults to `"dev"`. Goreleaser injects the actual version at build time via ldflags from the git tag — `version.go` never needs manual editing for releases.

Local dev builds (`task go-build`) embed the short git hash, UTC date, and `BuiltBy=local`.

### GitHub Actions

- **`main.yml`** — runs on every push; builds a snapshot with goreleaser to verify the build.
- **`release.yml`** — runs when a `v*` tag is pushed; performs the full goreleaser release.

## File Headers

All source files use SPDX license headers:
```go
// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT
```
