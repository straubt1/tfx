# TFx TUI — Planning Document

> Status: Phase 1 complete.
> Last updated: 2026-03-05 (Phase 1 implemented)

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
[Profile List]  ← future only; requires profile system
└── Organization  ← MVP: resolved from config, shown in header (not navigable)
    ├── Projects
    │   └── Workspaces
    │       ├── Runs
    │       │   └── Plan / Apply details
    │       ├── Configuration Versions
    │       ├── State Versions
    │       └── Variables
    ├── Variable Sets
    ├── Teams
    │   └── Team Members
    └── (future) Registry Modules, Registry Providers
```

### MVP: Config-resolved Organization

In the MVP, the organization is read from the existing config/env mechanism (Viper: `tfeOrganization` / `TFE_ORGANIZATION`), exactly as the CLI does today. The org is displayed in the header but is **not** a navigable level — the TUI drops directly into the project list.

```
tfx --tui
# → reads org from config/env → opens project list for that org
```

This requires no changes to how credentials work. The TUI entry point (`tui.Run()`) calls `client.NewFromViper()` and reads `viper.GetString("tfeOrganization")` just as every CLI command does.

### Future: Profile System

A profile system will allow users to configure multiple named profiles (hostname + org + token combinations), enabling the TUI to start at the **organization list level** and let the user select which org to enter.

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

With profiles configured, `tfx --tui` (no explicit org) would open a **profile/org picker** as the entry screen. If a `--tfeOrganization` flag or env var is set, it bypasses the picker and drops straight into the project list (consistent with current CLI behavior).

Navigation with profiles:
```
Profile/Org Picker  (future)
└── Projects        (MVP entry point)
    └── Workspaces
        └── ...
```

The breadcrumb bar reflects whichever level is active:
```
# MVP (org from config)
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

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Navigate rows |
| `Enter` | Drill into selected resource |
| `Esc` | Go up one level |
| `r` | Refresh current view |
| `/` | Filter / search |
| `?` | Show keybinding help overlay |
| `q` | Quit TUI |
| `c` | Copy CLI command hint to clipboard |
| `g` | Jump to top |
| `G` | Jump to bottom |

Additional context-specific shortcuts (e.g., `l` for logs on a Run) defined per view.

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

**In MVP:**
- [ ] `tfx --tui` flag wired into root command
- [ ] App shell: header, breadcrumb, status bar, CLI hint bar
- [ ] Workspace list view (projects → workspaces as the entry point)
- [ ] Run list view (drill into workspace → runs)
- [ ] Keyboard navigation (up/down/enter/esc/q/r/?)
- [ ] Filter (`/`) on workspace list
- [ ] CLI hint bar updates per view
- [ ] Terminal resize handling
- [ ] Basic Lip Gloss styling (header, selected row, status bar)
- [ ] Error display in status bar

**Post-MVP (future iterations):**
- [ ] Profile system (named profiles in `.tfx.hcl` with hostname + org + token)
- [ ] Profile/org picker as TUI entry point (when no org in config)
- [ ] Organization list view (list orgs accessible to the token)
- [ ] State versions view
- [ ] Variables view
- [ ] Variable sets view
- [ ] Teams / users view
- [ ] Configuration versions view
- [ ] Split panel (list + detail pane)
- [ ] Run detail / plan output viewport
- [ ] Clipboard copy of CLI hint (`c` key)
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

### Phase 3 — Run List View (second level)
1. Create `tui/views/runs.go`
2. Wire `data.FetchRuns()` with workspace context
3. Breadcrumb updates to show `org > workspace`
4. Esc returns to workspace list
5. CLI hint: `tfx workspace run list -n <workspace>`

### Phase 4 — Polish & Release
1. Refine styling (colors, borders, spacing)
2. Keyboard shortcut help overlay (`?`)
3. Ensure all terminal sizes work (min width/height guards)
4. README entry and `tfx --help` description
5. Integration test for TUI launch (basic smoke test)

---

## 11. Open Questions

- [x] **Resolved:** `tfx tui` subcommand. Flag-on-rootCmd approach was rejected because cobra's required-flag check runs before `PersistentPreRun`, preventing viper from satisfying the required org/token flags in time. Subcommand goes through `postInitCommands`/`presetRequiredFlags` normally.
- [ ] MVP entry point is project list (not workspace list) — confirm this is right. Workspace list is one more click but more directly useful day-to-day.
- [ ] `bubble-table` (Evertras) vs core Bubbles `table`? Evaluate during Phase 2 — core table may be sufficient.
- [ ] Profile system config format: HCL `profile` blocks in `.tfx.hcl` (see section 4)? Or a separate `~/.tfx-profiles.hcl`? To decide when we get there.
- [ ] When profile system lands, should `tfx --tui` with no org configured show the profile picker or error? Likely picker — but org-required flag on rootCmd will need to be relaxed for TUI mode.
- [ ] Clipboard support for CLI hint: `golang.design/x/clipboard` or `atotto/clipboard` — cross-platform, evaluate during Phase 4.

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
