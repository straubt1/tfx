# Development Setup

This guide covers setting up a macOS workstation for TFx development from scratch.

## Prerequisites

Install [Homebrew](https://brew.sh) if you don't have it:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

## Install Dependencies

Install the core tools:

```bash
brew install go goreleaser go-task
```

| Tool | Purpose |
|---|---|
| `go` | Go compiler and toolchain |
| `goreleaser` | Cross-platform builds and releases |
| `go-task` | Task runner (provides the `task` command) |

### Optional

```bash
brew install node   # only needed for serving the docs site locally (task development:serve-docs)
brew install --cask font-roboto-mono-nerd-font  # required for Nerd Font icons in VHS recordings
```

Verify everything is installed:

```bash
task doctor
```

## Clone and Build

```bash
git clone https://github.com/straubt1/tfx.git
cd tfx
go mod download
task go:build
./tfx version
```

## Common Tasks

Run `task --list` to see all available tasks:

| Task | Description |
|---|---|
| `task development:doctor` | Verify all required development tools are installed |
| `task development:serve-docs` | Serve documentation site locally (Starlight/Astro) |
| `task go:build` | Build Go binary (development build with git metadata) |
| `task go:build-all` | Cross-platform snapshot build via goreleaser (same as release:dry-run) |
| `task go:upgrade` | Upgrade Go toolchain (via brew) and all module dependencies |
| `task release:dry-run` | Simulate the full release pipeline locally (no publish, no tag required) |
| `task release:major` | Cut a major release (X+1.0.0) — for breaking changes |
| `task release:minor` | Cut a minor release (x.Y+1.0) — for new features |
| `task release:patch` | Cut a patch release (x.y.Z+1) — for bugfixes |
| `task test:integration-cmd` | Run integration tests for commands only |
| `task test:integration-data` | Run integration tests (requires TFE_HOSTNAME, TFE_TOKEN, TFE_ORGANIZATION env vars) |

## Upgrading Dependencies

To upgrade Go, goreleaser, and all module dependencies:

```bash
task go:upgrade
```

This will:
1. Upgrade `go` and `goreleaser` via Homebrew
2. Update `go.mod` to the installed Go version
3. Upgrade all module dependencies to their latest versions
4. Run `go mod tidy` to clean up

## Configuration for CLI / Integration Tests

TFx requires a TFE/HCP Terraform instance for integration tests. Create `secrets/.env-int`:

```bash
mkdir -p secrets
cat > secrets/.env-int << 'EOF'
TFE_HOSTNAME=app.terraform.io
TFE_TOKEN=your-api-token
TFE_ORGANIZATION=your-org
EOF
```

Then run integration tests:

```bash
task test:integration-data
task test:integration-cmd
```

## Project Structure

```
cmd/          # Cobra command handlers
  flags/      # Per-command flag structs
  views/      # Output/rendering for each command
client/       # TFE API client wrapper
data/         # Data fetching layer (API calls + business logic)
output/       # Output system (tables, JSON, spinner, logger)
tui/          # Bubble Tea TUI
integration/  # Integration tests
pkg/file/     # File utilities
version/      # Version info (injected at build time)
```

## Releasing

See the release tasks:

```bash
task release:patch   # bugfixes (x.y.Z+1)
task release:minor   # new features (x.Y+1.0)
task release:major   # breaking changes (X+1.0.0)
```

Always update `CHANGELOG.md` before cutting a release.
