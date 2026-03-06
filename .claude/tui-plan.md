# TFx TUI — Planning Document

> Status: Phases 1–6 complete.
> Last updated: 2026-03-06

---

## 1. Overview & Goals

Add a terminal user interface (TUI) to TFx, triggered via `tfx tui` (opt-in subcommand, leaving all existing CLI commands untouched). Inspired by k9s. Designed to be modern, keyboard-driven, and deeply integrated with TFE/HCP Terraform resource hierarchies.

**Non-goals for now:**
- Replacing the CLI — TUI is additive
- Mouse click support (planned for the future)
- Editing resources (read-only MVP)

---

## 2. Tech Stack Decision

### Selected: Bubble Tea v2 + Lip Gloss v2 + Bubbles v2

**Why Bubble Tea over tview:**

| Concern | Bubble Tea v2 | tview |
|---|---|---|
| Maintenance | Corporate-backed Charmbracelet team, v2 stable released Feb 2026 | Single maintainer, 44 open PRs, k9s had to fork it |
| Architecture | ELM / functional (testable, scales well) | Mutable widget tree (gets tangled at scale) |
| Styling | Lip Gloss — CSS-like, composable, beautiful | Per-widget color methods, no theme system |
| Resize handling | Manual (handle `WindowSizeMsg`) | Automatic in layout primitives |
| Layout | Manual math (weakness) | Grid/Flex auto-reflow (strength) |
| Long-term bet | Yes — 40k stars, growing ecosystem | Risky — bus factor, GitLab migrating away |

**tview is** what k9s uses and is the easier short-term path, but k9s had to fork it to get the production behavior they needed. Starting on Bubble Tea means we don't inherit that risk and we get a much better styling story, which matters for "modern and pretty."

**Resize weakness mitigation:** Bubble Tea's manual layout math is manageable if we design the layout system up front and centralize all size computation in one place.

**Key dependencies:**
```
charm.land/bubbletea/v2          # core TUI framework
charm.land/lipgloss/v2           # styling
charm.land/x/bubbles/v2          # built-in components (table, list, viewport, textinput, help, key)
github.com/Evertras/bubble-table # richer table with filtering, sorting (evaluate vs built-in)
```

> Note: Bubble Tea v2 changed its import path to `charm.land/...`. This is a vanity domain, not a GitHub path.

---

## 3. Entry Point Design

### Flag Approach
Add `--tui` as a persistent flag on the root command (or a dedicated `tui` subcommand). The simplest path is a boolean flag on `rootCmd`:

```
tfx --tui
```

When the flag is set:
1. Skip normal cobra command dispatch
2. Initialize TUI with the same Viper config (hostname, org, token) that CLI uses
3. Launch Bubble Tea program

The `tfeOrganization` and `tfeToken` flags remain required — the TUI reads credentials the same way the CLI does.

**Future option:** If we later decide TUI should be the default, we change the `rootCmd.RunE` to launch the TUI when called with no subcommand. The flag approach now gives us that migration path cleanly.

### Cobra integration sketch
```go
// cmd/root.go addition
var tuiMode bool
rootCmd.PersistentFlags().BoolVar(&tuiMode, "tui", false, "Launch the interactive TUI.")

// in Execute() or rootCmd.RunE:
if tuiMode {
    return tui.Run()
}
```

---

## 4. Resource Hierarchy / Navigation Model

Mirrors the TFE/HCP Terraform object model. Navigation is hierarchical (drill-down):

```
[Profile List]      ← future only; requires profile system
└── Organizations   ← Phase 6: navigable list of orgs accessible to the token
    └── Projects    ← MVP entry point (org resolved from config, not yet a list)
        └── Workspaces
            ├── Runs             ← enter (default drill-in from workspace list)
            │   └── Plan / Apply details
            ├── Variables        ← v key from workspace list
            ├── Config Versions  ← f key from workspace list ("files")
            └── State Versions   ← s key from workspace list
[Org-level views — future]
    ├── Variable Sets
    ├── Teams
    │   └── Team Members
    └── Registry Modules, Registry Providers
```

### Workspace navigation keys

From the workspace list, multiple shortcut keys drill into different workspace sub-views. This avoids a sub-menu screen and keeps navigation direct (k9s-style):

| Key | Destination |
|---|---|
| `enter` | Runs (most common — the default drill-in) |
| `v` | Variables |
| `f` | Configuration Versions ("files") |
| `s` | State Versions |

### Current: Config-resolved Organization (Phases 1–4)

In the current implementation, the organization is read from the existing config/env mechanism (Viper: `tfeOrganization` / `TFE_ORGANIZATION`), exactly as the CLI does today. The org is displayed in the header but is **not** a navigable level — the TUI drops directly into the project list.

```
tfx tui
# → reads org from config/env → opens project list for that org
```

This requires no changes to how credentials work. The TUI entry point (`tui.Run()`) calls `client.NewFromViper()` and reads `viper.GetString("tfeOrganization")` just as every CLI command does.

### Phase 6: Organization List View

Add organizations as a navigable top-level construct. The TUI will call `client.Organizations.List()` to enumerate all orgs accessible to the configured token, then present them as a selectable list. Selecting an org drills into its project list.

**Behavior:**
- If `TFE_ORGANIZATION` (or equivalent) is configured, the TUI still opens at the org list but pre-selects / highlights that org — giving the user the choice to enter or switch.
- If no org is configured, the org list is the mandatory entry screen.
- `Enter` on an org selects it and navigates to the project list for that org.
- `Esc` from the project list returns to the org list.

**Breadcrumb evolution:**
```
# Phase 6 (org list as entry)
org: my-org  >  project: platform-team  >  workspace: prod-app

# With org list visible
[orgs]  >  my-org  >  project: platform-team  >  workspace: prod-app
```

**Implementation sketch:**
- `viewOrganizations` added to `viewType` enum
- `orgs []*tfe.Organization`, `orgCursor`, `orgOffset`, `orgFilter`, `orgFiltering` added to `Model`
- `orgsLoadedMsg` and `loadOrganizations(c)` in `messages.go`
- `tui/organizations.go` — `orgColumns()`, `filteredOrgs()`, `renderOrgsContent()`
- `handleOrgsKey()` — `enter` → set `m.selectedOrg`, `loadProjects` for that org, transition to `viewProjects`
- Initial `Init()` dispatches `loadOrganizations` instead of `loadProjects`
- CLI hint: `tfx organization list` at org-list level

**Data function needed:**
`data.FetchOrganizations(c *client.TfxClient) ([]*tfe.Organization, error)` — wraps `c.Client.Organizations.List()` with pagination. Check if this already exists in `data/` before creating.

### Future: Profile System

A profile system will allow users to configure multiple named profiles (hostname + org + token combinations) as a higher level above the org list.

Profiles would live in the config file (`.tfx.hcl`) under a `profile` block:

```hcl
profile "prod" {
  hostname     = "app.terraform.io"
  organization = "acme-prod"
  token        = "..."
}

profile "staging" {
  hostname     = "tfe.internal.acme.com"
  organization = "acme-staging"
  token        = "..."
}
```

With profiles configured, `tfx tui` (no explicit org) would open a **profile picker** as the entry screen, then the org list for the selected profile's token.

Navigation with profiles:
```
Profile Picker  (future)
└── Organizations  (Phase 6 entry point)
    └── Projects
        └── Workspaces
            └── ...
```

The breadcrumb bar reflects whichever level is active:
```
# Phase 6 (org list as entry)
org: my-org  >  project: platform-team  >  workspace: prod-app

# Future (with profiles)
profile: prod  >  org: my-org  >  project: platform-team  >  workspace: prod-app
```

Navigation pattern: Enter drills in, Esc goes up one level.

---

## 5. Layout Design

### Zones (fixed across all views)

```
┌─────────────────────────────────────────────────────────┐
│ HEADER: hostname · org · version             [keybinds] │  ← 1 line
├─────────────────────────────────────────────────────────┤
│ BREADCRUMB: org > project > workspace                   │  ← 1 line
├─────────────────────────────────────────────────────────┤
│                                                         │
│  MAIN CONTENT AREA                                      │  ← fills remainder
│  (table, detail pane, etc.)                             │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ STATUS BAR: loading... / error msg                      │  ← 1 line
├─────────────────────────────────────────────────────────┤
│ CLI HINT: tfx workspace show -n prod-app                │  ← 1 line
└─────────────────────────────────────────────────────────┘
```

> **Why does the header show `vdev`?**
> The header renders `v{version.Version}`. The `version.Version` variable (in `version/version.go`) defaults to the string `"dev"` in source. Goreleaser injects the real git tag at release time via ldflags (`-X github.com/straubt1/tfx/version.Version=v1.2.3`). In local development builds (`task go-build`), goreleaser is not invoked, so the version stays as `"dev"` and the header shows `vdev`. This is expected and correct — no action needed.

**CLI hint bar** (bottom): always shows the equivalent `tfx` CLI command for the current view. This is a core UX requirement — users who want to script or automate can easily discover the right command.

### Content area layouts (per view type)

**List view** (default for most resources):
```
┌─────────────────────────────────────────────────────────┐
│ [/] filter...                                           │
├────────────────────────────────────────────────────────-┤
│ NAME          STATUS    UPDATED         TERRAFORM VER   │
│ > prod-app    Active    2h ago          1.9.5            │
│   staging     Active    1d ago          1.9.3            │
│   dev         Locked    3d ago          1.8.0            │
└─────────────────────────────────────────────────────────┘
```

**Split view** (detail panels — future, post-MVP):
- Left: list (30-40% width)
- Right: detail/YAML pane (60-70% width)

---

## 6. Keyboard Shortcuts

**Global (all views):**

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Navigate rows |
| `Esc` | Go up one level / clear filter |
| `r` | Refresh current view |
| `/` | Filter / search |
| `?` | Show keybinding help overlay |
| `q` | Quit TUI |
| `c` | Copy CLI command hint to clipboard |
| `g` | Jump to top |
| `G` | Jump to bottom |

**Workspace list (context-specific):**

| Key | Action |
|---|---|
| `enter` | Drill into runs (default) |
| `v` | Drill into variables |
| `f` | Drill into configuration versions |
| `s` | Drill into state versions |

Additional context-specific shortcuts (e.g., `a` to apply a planned run) defined per view.

Keybindings will use the Bubbles `key` package to define and auto-render help.

---

## 7. Styling Direction

"Modern and pretty" — draw inspiration from k9s's color density, but with Charm's cleaner aesthetic.

**Palette (proposed — to be iterated):**
- Background: terminal default (transparent)
- Header bg: deep purple / navy
- Header fg: white
- Selected row: bright highlight (cyan or orange accent)
- Dimmed text: gray
- Error: red
- Success / active: green
- CLI hint bar: italic, subdued

Lip Gloss `AdaptiveColor` for light/dark terminal auto-adaptation.

### Loading Animation

The current "Loading…" static text in the status bar is boring. Replace it with an animated spinner using the **`charm.land/x/bubbles/v2` spinner component** — no new external dependency, since Bubbles v2 is already in the tech stack.

**How the Bubbles spinner works (BT v2):**

The `bubbles/spinner` package provides a `spinner.Model` that holds a frame index and a `Tick` command. Each `Tick` produces a `spinner.TickMsg` that advances the frame. In your `Update()`:

```go
case spinner.TickMsg:
    if m.loading {
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
```

The spinner is only ticked while loading is in progress; once data arrives, ticking stops naturally (no more `Tick` commands are dispatched).

**Built-in spinner styles (choose one):**

| Name | Frames | Look |
|---|---|---|
| `spinner.Dot` | ⣾⣽⣻⢿⡿⣟⣯⣷ | Braille dot sweep (default) |
| `spinner.Line` | \|/-\ | Classic ASCII line |
| `spinner.MiniDot` | ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ | Smaller braille |
| `spinner.Jump` | ⢄⢂⢁⡁⡈⡐⡠ | Jumping dot |
| `spinner.Pulse` | █▓▒░ | Pulsing blocks |
| `spinner.Points` | ∙∙∙ / •∙∙ etc | Three dots |
| `spinner.Globe` | 🌍🌎🌏 | Rotating globe emoji |
| `spinner.Moon` | 🌑🌒🌓🌔🌕🌖🌗🌘 | Moon phases emoji |

**Recommended:** `spinner.Dot` (braille sweep) or `spinner.MiniDot` — looks modern, no emoji dependency, works in all terminals.

**Implementation plan:**

1. Add `charm.land/x/bubbles/v2` to `go.mod` (run `go get charm.land/x/bubbles/v2`)
2. Add `spinner spinner.Model` to `Model` struct; initialize in `newModel()` with preferred style and accent color
3. In `Init()`, return `m.spinner.Tick` alongside any initial fetch command
4. In `Update()`, handle `spinner.TickMsg` → only propagate tick when `m.loading == true`
5. In the status bar render, replace `"Loading…"` with `m.spinner.View() + " Loading…"` when loading

---

## 8. Code Structure

New top-level package `tui/`:

```
tui/
  model.go          # root Bubble Tea model, layout, WindowSizeMsg handling
  keys.go           # global keybindings
  styles.go         # Lip Gloss style definitions (palette, reusable styles)
  header.go         # header component (hostname, org, version)
  breadcrumb.go     # breadcrumb bar component
  statusbar.go      # status / error bar
  clihint.go        # CLI command hint bar
  views/
    projects.go     # project list view
    workspaces.go   # workspace list view
    runs.go         # run list view
    variables.go    # variable list view
    varsets.go      # variable set list view
    teams.go        # team list view
    detail.go       # generic detail / viewport view
```

**Data layer:** TUI views call the same `data/` functions the CLI uses. No new data fetching code — purely a new rendering layer. The `client.TfxClient` is initialized once from Viper (same as CLI) and threaded through TUI commands.

---

## 9. MVP Scope

MVP = enough to demo the concept and gather feedback. Scope:

**Delivered (Phases 1–4):** ✅
- [x] `tfx tui` subcommand wired into cobra
- [x] App shell: header, breadcrumb, status bar, CLI hint bar
- [x] Project list view (entry point)
- [x] Workspace list view (drill from project)
- [x] Run list view (drill from workspace via `enter`) with colored status
- [x] Keyboard navigation (up/down/enter/esc/q/r/g/G/?)
- [x] Filter (`/`) on all list views
- [x] CLI hint bar updates per view
- [x] `c` key copies CLI hint to clipboard
- [x] Terminal resize handling + min-size guard
- [x] GitHub Dark palette styling
- [x] Loading and error states in status bar

**Phase 5 — Workspace sub-views:** ✅
- [x] Variables list (`v` key from workspace list) — `tfx workspace variable list -n <ws>`
- [x] Configuration Versions list (`f` key) — `tfx workspace cv list -n <ws>`
- [x] State Versions list (`s` key) — `tfx workspace sv list -n <ws>`
- [x] Breadcrumb and CLI hint update for each sub-view
- [x] Esc from any sub-view returns to workspace list

**Phase 5.5 — Loading Animation:** ✅
- [x] Manual spinner (no extra dependency) — `spinnerFrames []string` braille sweep, driven by `spinnerTickMsg` / `tickSpinner()` command chain
- [x] Spinner animates in status bar and loading content area while `m.loading == true`
- [x] Stops cleanly when data arrives (chain terminates when `m.loading` is false)

**Phase 6 — Organization List View:** ✅
- [x] `viewOrganizations` added to `viewType` enum (entry point iota = 0)
- [x] Org state fields added to `Model`
- [x] `orgsLoadedMsg` and `loadOrganizations(c)` in `messages.go`
- [x] `tui/organizations.go` with `orgColumns()`, `filteredOrgs()`, `renderOrgsContent()`
- [x] `Init()` dispatches `loadOrganizations` + `tickSpinner()` via `tea.Batch`
- [x] `handleOrgsKey()`: enter selects org, sets `m.org`, triggers `loadProjects`
- [x] Breadcrumb: `organizations` active at org level
- [x] CLI hint: `tfx organization list` at org list level
- [x] Pre-highlight configured org on load (`TFE_ORGANIZATION` match)
- [x] `Esc` from project list returns to org list (no re-fetch)

**Post-MVP (future iterations):**
- [ ] Profile system (named profiles in `.tfx.hcl` with hostname + org + token)
- [ ] Profile picker as TUI entry point (level above org list)
- [ ] Variable sets view (org-level)
- [ ] Teams / users view
- [ ] Split panel (list + detail pane)
- [ ] Run detail / plan output viewport
- [ ] Run actions (queue run, cancel run — write operations)
- [ ] Mouse click support
- [ ] Custom theme / color scheme config
- [ ] Make TUI the default `tfx` behavior (no flag needed)
- [ ] Registry modules / providers view

---

## 10. Implementation Plan (phased)

### Phase 1 — Shell & Wiring (no data, static UI) ✅ COMPLETE
1. ~~Add `--tui` flag to `cmd/root.go`~~ → `tfx tui` subcommand in `cmd/tui.go`
2. Created `tui/` package: `run.go`, `model.go`, `styles.go`
3. Layout zones implemented: header, breadcrumb, content (placeholder), status bar, CLI hint
4. `WindowSizeMsg` handling and layout math in place
5. Global keybindings: `q`/`ctrl+c` quit, `?` help overlay, `esc` close help
6. Lip Gloss v2 styles in `tui/styles.go` (GitHub Dark palette)
7. `tui.Run()` wired into `cmd/tui.go` cobra subcommand
8. Build passes, all tests pass

**Bubble Tea v2 API notes (for future phases):**
- `Init() tea.Cmd` (not `(Model, Cmd)` — research was wrong)
- `View() tea.View` (not `string` like v1) — use `tea.NewView(content)`
- Alt screen: `view.AltScreen = true` on the returned View (not a `ProgramOption`)
- Key events: `tea.KeyPressMsg` (not `tea.KeyMsg` like v1)
- `tea.Quit` is still a `func() Msg` = `Cmd`, usage unchanged

### Phase 2 — Workspace List View (first live data)
1. Create `tui/views/workspaces.go` using Bubbles `table` or `bubble-table`
2. Wire `data.FetchWorkspaces()` call (async via `tea.Cmd`)
3. Implement loading spinner in status bar during fetch
4. Implement error display in status bar
5. Implement `/` filter
6. CLI hint bar shows `tfx workspace list` / `tfx workspace show -n <name>`
7. Navigate rows, press Enter → placeholder "drill" message

### Phase 3 — Run List View (second level) ✅ COMPLETE
1. Created `tui/runs.go` with run list, colored status, relative timestamps
2. Wired `data.FetchRunsForWorkspace()` (max 50 runs) via async `tea.Cmd`
3. Breadcrumb: `org / project: X / workspace: Y / runs`
4. Esc returns to workspace list; `r` refreshes
5. CLI hint: `tfx workspace run list -n <workspace>`

### Phase 4 — Polish & Release ✅ COMPLETE
1. Run status colors: green=applied, amber=in-progress, blue=planned, red=errored, dim=canceled
2. `c` key copies current CLI command to clipboard (pbcopy/xclip/clip); green status bar feedback
3. `renderTableRowWithCellStyles` for per-cell foreground overrides
4. `colorSuccess` / `statusSuccessStyle` added to palette
5. Help overlay updated with all current keybindings

### Phase 5 — Workspace Sub-Views
Goal: from the workspace list, navigate directly into Variables, Config Versions, or State Versions via shortcut keys.

**Navigation keys from workspace list:**
- `v` → Variables list for selected workspace
- `f` → Configuration Versions list
- `s` → State Versions list
- `enter` → Runs (existing, unchanged)

**New files needed:**
- `tui/variables.go` — `renderVariablesContent()`, `filteredVariables()`, `variableColumns()`
- `tui/configversions.go` — `renderConfigVersionsContent()`, `filteredConfigVersions()`, `cvColumns()`
- `tui/stateversions.go` — `renderStateVersionsContent()`, `filteredStateVersions()`, `svColumns()`

**model.go changes:**
- Add `viewVariables`, `viewConfigVersions`, `viewStateVersions` to the `viewType` enum
- Add state fields for each (data slice, cursor, offset, filter, filtering flag)
- Add `v`, `f`, `s` key handling in `handleWorkspacesKey()`
- Update `navigateBack()`, `refresh()`, `isFiltering()`, `handleFilterKey()`
- Update breadcrumb, status bar, CLI hint for each new view

**messages.go changes:**
- Add `variablesLoadedMsg`, `configVersionsLoadedMsg`, `stateVersionsLoadedMsg`
- Add `loadVariables()`, `loadConfigVersions()`, `loadStateVersions()` fetch commands

**Data functions to use:**
- `data.FetchVariables(c, orgName, workspaceName)` for Variables
- `data.FetchConfigurationVersions(c, workspaceID)` for Config Versions
- `data.FetchStateVersions(c, workspaceID)` for State Versions

### Phase 5.5 — Animated Loading Spinner
Goal: replace the static "Loading…" text with a Bubble Tea–native animated spinner.

1. `go get charm.land/x/bubbles/v2` to pull in the spinner package (it may already be indirect)
2. Add `spinner spinner.Model` field to `Model`; initialize in `newModel()`:
   ```go
   s := spinner.New()
   s.Spinner = spinner.Dot
   s.Style = lipgloss.NewStyle().Foreground(colorAccent)
   m.spinner = s
   ```
3. In `Init()`, add `m.spinner.Tick` to the batch of initial commands
4. In `Update()`, add a case for `spinner.TickMsg`:
   ```go
   case spinner.TickMsg:
       if m.loading {
           m.spinner, cmd = m.spinner.Update(msg)
           cmds = append(cmds, cmd)
       }
   ```
5. In `renderStatusBar()`, replace `"Loading…"` with `m.spinner.View() + " Loading…"` when `m.loading`
6. Spinner stops automatically when loading transitions to false (no more Tick dispatched)

### Phase 6 — Organization List View
Goal: add organizations as a top-level navigable construct so users can switch orgs without restarting.

1. Add `data.FetchOrganizations(c *client.TfxClient) ([]*tfe.Organization, error)` to `data/` (check if it already exists)
2. Add `orgsLoadedMsg []*tfe.Organization` and `loadOrganizations(c)` cmd to `messages.go`
3. Create `tui/organizations.go`:
   - `orgColumns(width int)` — columns: Name, Email, Created, Plan
   - `filteredOrgs(m Model) []*tfe.Organization`
   - `renderOrgsContent() string`
4. Add `viewOrganizations` to the `viewType` enum (insert before `viewProjects` — it's the new root)
5. Add org state fields to `Model`:
   ```go
   orgs        []*tfe.Organization
   orgCursor   int
   orgOffset   int
   orgFilter   string
   orgFiltering bool
   selectedOrg *tfe.Organization
   ```
6. In `newModel()`, set initial view to `viewOrganizations` and dispatch `loadOrganizations`
7. Add `handleOrgsKey()` method:
   - `enter` → set `m.selectedOrg`, clear projects, dispatch `loadProjects(c, org.Name)`, transition to `viewProjects`
   - `/` → start filtering; `esc` → clear filter / return to org list (no level above)
   - `q` → quit
8. Update `navigateBack()`: from `viewProjects` → `viewOrganizations` (reset project list)
9. Update `refresh()`: when `viewOrganizations`, dispatch `loadOrganizations`
10. Update breadcrumb: show `[orgs]` at org list level; show org name when inside
11. Update CLI hint: `tfx organization list` at org list level
12. If `TFE_ORGANIZATION` env/config is set, highlight that row in the list on load (user can press enter to confirm or pick a different one)

---

## 11. Open Questions

- [x] **Resolved:** `tfx tui` subcommand (not `--tui` flag). Cobra's required-flag check runs before `PersistentPreRun`; subcommand goes through `postInitCommands`/`presetRequiredFlags` normally.
- [x] **Resolved:** Entry point is project list → workspaces (confirmed correct; org from config).
- [x] **Resolved:** Hand-rolled table renderer (no external library). `bubble-table` is BT v1 only; `bubbles/table` v2 lacks filtering and row-metadata association needed for drill-in navigation.
- [x] **Resolved:** Clipboard via platform-native exec (pbcopy/xclip/clip) — no external dependency needed.
- [ ] **Phase 6:** Does `data.FetchOrganizations` already exist? Check `data/` before implementing. If the token only has access to one org (common in single-org TFE installs), consider auto-advancing past the org list.
- [ ] **Phase 6:** When `TFE_ORGANIZATION` is set and Phase 6 lands, should the TUI skip the org list entirely and go straight to projects (matching current CLI behavior), or always show the org list with the configured org pre-highlighted? Pre-highlighted is friendlier but slightly slower (requires an API call). Leaning toward pre-highlighted.
- [ ] Profile system config format: HCL `profile` blocks in `.tfx.hcl` (see section 4)? Or a separate `~/.tfx-profiles.hcl`? To decide when we get there.
- [ ] When profile system lands, should `tfx tui` with no org configured show the profile picker or error? Likely picker — but org-required flag on rootCmd will need to be relaxed for TUI mode.

---

## 12. References

- [Bubble Tea v2](https://charm.land/bubbletea/v2) — Elm architecture TUI framework
- [Lip Gloss v2](https://charm.land/lipgloss/v2) — CSS-like terminal styling
- [Bubbles v2](https://charm.land/x/bubbles/v2) — component library (table, list, viewport, key, help)
- [bubble-table by Evertras](https://github.com/Evertras/bubble-table) — extended table component
- [k9s](https://k9scli.io/) — UX inspiration
- [Bubble Tea layout tips](https://leg100.github.io/en/posts/building-bubbletea-programs/)
- [gh-dash](https://github.com/dlvhdr/gh-dash) — Bubble Tea reference implementation (similar concept: GitHub resource browser)
- TFE/HCP Terraform API: `github.com/hashicorp/go-tfe`
