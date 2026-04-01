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

**Do not use worktrees.** Never use `isolation: "worktree"` or any equivalent when spawning agents or sub-tasks.

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
- `--config-file` flag or `TFX_CONFIG_FILE` env var — explicit config file path; auto-discovered at `./.tfx.hcl` (current dir) or `~/.tfx.hcl` (home dir) when neither is set

Integration tests load secrets from `secrets/.env-int` (not in repo).

## Project Structure

```
cmd/          # Cobra command implementations and orchestration
  flags/      # Per-command flag structs and parsers
  views/      # Output/rendering for each command
client/       # TFE API client wrapper (wraps go-tfe)
data/         # Data fetching layer — API calls and business logic
output/       # Output system (terminal tables, JSON, spinner, logger)
tui/          # Bubble Tea TUI — models, renderers, styles
  run.go      # Entry point (tui.Run())
  model.go    # Root model, key dispatch, layout
  styles.go   # Global lipgloss style vars
  debugpanel.go # API inspector panel
  # one file per view/feature
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

## TUI Architecture & Patterns

> Visual layout reference (ASCII diagrams for all views): `tui/LAYOUT.md`

The TUI lives in `tui/` and uses Bubble Tea v2 + Lip Gloss v2. It shares the same `data/` and `client/` layers as the CLI. Entry point: `cmd/tui.go` → `tui.Run()`.

### Key TUI dependencies
- `charm.land/bubbletea/v2` — ELM architecture event loop
- `charm.land/lipgloss/v2` — terminal styling
- `github.com/charmbracelet/x/ansi` — ANSI-safe string truncation (transitive dep of lipgloss v2, no extra `go get` needed)

### Bubble Tea v2 API — common gotchas
These differ from the widely-documented v1 API and will cause compile errors:
- `Init() tea.Cmd` — NOT `(Model, Cmd)`
- `View() tea.View` — return `tea.NewView(content)`, NOT a plain string; set `view.AltScreen = true` on the returned View (there is no `tea.WithAltScreen()` option)
- Key events are `tea.KeyPressMsg`, NOT `tea.KeyMsg`

### Lip Gloss v2 — common gotchas
- `lipgloss.Color` is a **function** (`func(s string) color.Color`), not a type. Use `color.Color` from `"image/color"` for struct fields and function parameters that store a colour value.
- Measure rendered string width with `lipgloss.Width(s)` (ANSI-aware).
- Right-align within a fixed-width column: `style.Width(n).Align(lipgloss.Right).Render(s)`

### ANSI-safe truncation
Use `ansi.Truncate(str, width, "")` from `github.com/charmbracelet/x/ansi` to clip styled strings without breaking escape sequences.

### Multi-panel layout: width enforcement
Any layout with side-by-side panels needs a single enforcing function that *both* pads short lines AND truncates overwide ones using `ansi.Truncate`. Without explicit truncation, content from one panel bleeds into the adjacent separator column.

### Style threading for focus-aware panels
When a panel changes background based on focus state, styles cannot be global package vars — they need to be dynamic. Define a styles struct and pass it through every renderer:

```go
type panelStyles struct {
    bg      lipgloss.Style
    punct   lipgloss.Style
    panelBg color.Color  // raw colour for helpers that build styles dynamically
}
```

Functions that render coloured glyphs (e.g. HTTP method verbs, status codes) need `bg color.Color` as a parameter so the glyph background matches the panel, not a hardcoded global.

### Key dispatch ordering
Three-tier dispatch prevents focus-escape bugs where global keys intercept input meant for a focused sub-panel:
1. Always-global (quit, toggle panel, switch focus)
2. Focused panel — guarded by `panelFocused && panelVisible`
3. Main view (filter input, globals, view-specific keys)

## Login TUI (`tfx login`)

The login flow is an **inline** Bubble Tea TUI (no alt-screen) implemented in `tui/login.go`. Entry point: `cmd/login.go` calls `output.Get().DisableSpinner()` then `tui.RunLogin(hostname)`.

### Profile Properties

| Property | HCL key | Default | Required |
|---|---|---|---|
| Name | block label | `"default"` | — |
| Hostname | `tfeHostname` | `"app.terraform.io"` | — |
| Organization | `tfeOrganization` | (commented placeholder) | — |
| Token | `tfeToken` | — | ✓ |

`tfe-license-path` is a valid key in the block but is **not** written by `tfx login` and is not populated in the `Profile` struct — it passes through untouched in the raw file.

### HCL Config Format (`pkg/hclconfig/`)

**New format** — one or more named profile blocks:

```hcl
profile "default" {
  tfeHostname     = "app.terraform.io"
  tfeOrganization = "my-org"
  tfeToken        = "abc123..."
}

profile "staging" {
  tfeHostname     = "tfe.myco.internal"
  tfeOrganization = "my-org"
  tfeToken        = "xyz..."
}
```

- Block label = profile **name** (a user-editable alias — NOT the hostname).
- `tfeHostname` inside the block = hostname. Always independent of the name.
- `tfeHostname` is optional; defaults to `DefaultHostname` (`"app.terraform.io"`) when absent.
- `tfeOrganization` is optional; when empty `WriteProfile` writes a commented placeholder:
  ```
    # tfeOrganization = "" # set this to your organization name
  ```

**Legacy flat format** (no `profile` blocks) — still supported for reading:

```hcl
tfeHostname     = "app.terraform.io"   # optional, defaults to app.terraform.io
tfeOrganization = "my-org"
tfeToken        = "abc123..."
```

Parsed as a single profile with `Name = DefaultProfileName` and `Hostname = tfeHostname` value (or `DefaultHostname` if absent). Only returned when `tfeToken` is non-empty.

**Key constants:**
- `hclconfig.DefaultProfileName = "default"`
- `hclconfig.DefaultHostname = "app.terraform.io"`

**Key functions:**
- `Profile` struct: `Name`, `Hostname`, `Organization`, `Token string`
- `ListProfiles(path string) ([]Profile, error)` — `nil, nil` when file not found or empty
- `WriteProfile(path, name, hostname, organization, token string) error` — name/hostname default to constants when empty

**Backward compat for old block-label-as-hostname files:**
If `tfeHostname` is absent AND the block label contains a dot (e.g. `profile "app.terraform.io" {}`), the label is used as hostname. Aliases without dots (`"default"`, `"prod"`) fall back to `DefaultHostname`. The profile name is **never** used as hostname for name-only aliases.

### Profile Resolution in `cmd/root.go`

`resolveActiveProfile()` runs in `PersistentPreRunE` (after `bindPFlags`):

1. Load profiles from `viper.ConfigFileUsed()` via `hclconfig.ListProfiles`
2. If `--profile` flag set → find profile by `Profile.Name`, promote values to Viper
3. If `--profile` not set → prefer a profile named `"default"`, then fall back to `profiles[0]`
4. Always calls `viper.Set("profile", active.Name)` so the TUI can read it back
5. Calls `viper.Set("tfeHostname", active.Hostname)` **only when `Hostname != ""`** (avoids overriding env/flag defaults with an empty value)
6. Always sets `tfeToken` and `tfeOrganization`
7. Hostname is **NEVER** derived from the profile name — always from `tfeHostname` inside the block

### State Machine

```
stepProfileList  — existing-profile selector (skipped when no profiles in ~/.tfx.hcl)
stepProfileName  — text input for profile name, pre-filled with "default"
stepMenu         — two options: open browser / enter token directly
stepToken        — masked token input (● per char); paste via tea.PasteMsg
stepValidating   — spinner while fetching orgs from the API
stepTokenError   — validation failed: re-enter / accept anyway
stepOrgSelect    — arrow-key org picker (only when 2+ orgs returned)
stepDone         — success; shows created vs updated message
stepError        — fatal write/config error
stepCancelled    — clean exit (q / esc / ctrl+c)
```

### Entry Logic in `RunLogin()`

- If `ListProfiles(configPath)` returns 1+ profiles → start at `stepProfileList`
- Otherwise (no file, empty, error) → start at `stepProfileName` with `nameRunes` pre-filled to `[]rune("default")`

### Step-by-Step Flow

**Adding a new profile** (cursor 0 on profile list, or no existing profiles):

1. `stepProfileName` — text input pre-filled with `"default"`:
   - User accepts as-is (Enter) or edits the name, then Enter to confirm → advance to `stepMenu`
   - `nameRunes` is initialized to `[]rune(hclconfig.DefaultProfileName)` so pressing Enter immediately uses `"default"`
2. `stepMenu` — choose auth method:
   - `"Open browser to create a token"` → opens `https://<hostname>/app/settings/tokens?source=tfx-login`, advance to `stepToken`
   - `"Enter token directly"` → advance to `stepToken`
3. `stepToken` — masked input (● per char); show `✓ looks right` when token contains `.atlasv1.`; Enter to validate
4. `stepValidating` — call TFE API `Organizations.List` with the entered token
   - API error / 0 orgs → `stepTokenError`
   - 1 org → `finalize()` with that org auto-selected
   - 2+ orgs → `stepOrgSelect`
5. `stepOrgSelect` — pick org with ↑/↓, Enter → `finalize()`
6. `stepDone` — shows "Profile for \<hostname\> has been **created**"

**Re-authenticating an existing profile** (cursor N≥1 on profile list):
- Sets `hostname = profiles[N-1].Hostname`, `selectedProfileName = profiles[N-1].Name`, `isUpdate = true`
- Skips `stepProfileName` — goes directly to `stepMenu`
- `stepToken` shows amber overwrite warning: `⚠ Re-authenticating <hostname> — this will replace the existing token.`
- `stepDone` shows "Profile for \<hostname\> has been **updated**"

**Token validation failure** (`stepTokenError`):
- Option 0 `"Re-enter token"`: clear `tokenRunes` and `resolvedToken`, go back to `stepToken`
- Option 1 `"Accept anyway"`: set `selectedOrg = ""`, call `finalize()` immediately (no org prompt)
  - `finalize()` writes profile with empty org → commented placeholder in HCL
  - `stepDone` shows amber warning: `⚠ Organization not set — edit <configPath> to configure`

**`finalize()`** — called from `stepOrgSelect`, single-org path, or "accept anyway":
- `name = selectedProfileName` (fallback to `"default"` if empty — never uses hostname)
- Calls `hclconfig.WriteProfile(configPath, name, hostname, selectedOrg, resolvedToken)`
- Sets `stepDone` on success, `stepError` on failure
- Returns `tea.Quit`

### ESC / Back Navigation

| Step              | ESC goes to                                             |
|-------------------|---------------------------------------------------------|
| `stepProfileList` | `stepCancelled` (quit)                                  |
| `stepProfileName` | `stepProfileList` if profiles exist, else `stepCancelled` |
| `stepMenu`        | `stepProfileName` (new profile) or `stepProfileList` (re-auth) |
| `stepToken`       | `stepMenu`                                              |
| `stepOrgSelect`   | `stepToken` (clears token for re-entry)                 |
| `stepTokenError`  | `stepToken` (keeps token for editing)                   |

### Model Fields

```go
type LoginModel struct {
    step                loginStep
    hostname            string              // target hostname
    configPath          string
    profiles            []hclconfig.Profile // existing profiles from config
    profileCursor       int                 // stepProfileList cursor
    isUpdate            bool                // true when re-authing existing profile
    selectedProfileName string              // name written to profile block
    nameRunes           []rune              // stepProfileName text buffer (pre-filled "default")
    menuCursor          int                 // stepMenu: 0=browser, 1=direct
    useBrowser          bool
    tokenRunes          []rune
    tokenErr            error               // validation error shown in stepTokenError
    tokenErrCursor      int                 // 0=re-enter, 1=accept anyway
    orgs                []*tfe.Organization
    orgCursor           int
    selectedOrg         string
    resolvedToken       string
    spinnerIdx          int
    err                 error               // fatal write/config error
    width               int
}
```

### Rendering

- Inline mode: `tea.NewView(content)` — **no** `view.AltScreen = true`
- Header bar: full-width `colorHeaderBg` strip with `colorAccent` title + `colorDim` hostname
- Divider: `strings.Repeat("─", w)` in `colorBorder`
- `renderWidth()`: caps at 88, defaults to 72 before first `tea.WindowSizeMsg`
- Masked token: `strings.Repeat("●", len(tokenRunes))`
- Spinner frames: `[]string{"⠋","⠙","⠹","⠸","⠼","⠴","⠦","⠧","⠇","⠏"}`, tick every 80 ms
- Colors: `colorAccent` (selected/cursor), `colorDim` (inactive), `colorSuccess` (✓), `colorError` (✗), `colorLoading` (amber warnings), `colorBorder` (divider)
- Hints: dim italic footer on every step, e.g. `"Enter to continue · Backspace to delete · Esc to go back"`

### Key Implementation Notes

- `tea.PasteMsg` — use `msg.Content` field (struct, not string); strip `\n \r \t` for token, strip `\n \r` for name
- `isPrintable(k string) bool` — shared helper in `tui/debugpanel.go` (same package); used in both token and name entry
- `finalize()` — calls `hclconfig.WriteProfile`, sets `stepDone` or `stepError`, returns `tea.Quit`
- `--profile` flag in `cmd/root.go` matches by `Profile.Name`; hostname always read from `Profile.Hostname`
- `loginFetchOrgs` constructs its own `*tfe.Client` directly (not `client.NewFromViper`) because credentials are not in Viper yet at login time

## TUI Profile Bar (`tui/model.go` — `renderProfileBar`)

The profile bar sits between the main header and breadcrumb. It renders **four fixed rows** so `fixedLines` stays constant (no layout shift when data arrives):

```
  profile:   default
  username:  tstraub
  email:     tstraub@hashi.com
  expires:   2027-01-15
```

Labels are padded to 9 chars so values align vertically. All rows use `colorHeaderBg` background; labels use `colorDim`, values use `colorAccent`. Placeholder `"…"` is shown in `colorDim` until data loads.

### Data sources

- **profile / hostname**: available immediately from Viper (`profileName` set via `resolveActiveProfile`)
- **username / email**: loaded by `loadAccount` → `Users.ReadCurrent` → `accountLoadedMsg`
- **expires**: loaded by `loadAccountToken` → `UserTokens.List(userID)` after user is known → `accountTokenLoadedMsg`
  - Picks the token with the most-recent `LastUsedAt` (proxy for "the token currently in use")
  - `ExpiredAt.IsZero()` → displays `"never"`; any error → displays `"n/a"`
  - `loadAccountToken` is fired from the `accountLoadedMsg` handler once the user ID is available

## Key Dependencies

- `github.com/hashicorp/go-tfe` — TFE/HCP Terraform API client
- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — Config and flag management
- `github.com/jedib0t/go-pretty` — Terminal table formatting
- `github.com/fatih/color` — Colored terminal output
- `github.com/pkg/errors` — Error wrapping

## Release Process

Releases are automated via goreleaser and GitHub Actions. The release workflow triggers on a new git tag and publishes binaries, a Docker image (GHCR), Linux packages (apk/deb/rpm), and a Homebrew formula update to `straubt1/homebrew-tap`.

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

Runs `goreleaser release --snapshot --clean --skip=announce,validate` — builds all artifacts without requiring a tag. The Homebrew formula is generated and written to `dist/` (not pushed to the tap) because `skip_upload: "{{ .IsSnapshot }}"` is set in `.goreleaser.yml`. Goreleaser v2 may emit 2 informational `dockers_v2` warnings; these are a known goreleaser quirk and can be ignored.

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
