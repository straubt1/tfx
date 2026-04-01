---
title: Navigation
description: View hierarchy, keyboard shortcuts, and navigation patterns in the TFx TUI.
---

## View hierarchy

The TUI is organized as a drill-down hierarchy. Press `enter` to go deeper, `esc` to go back. A breadcrumb bar at the top of the screen always shows your current location.

```
Organizations
  └── Projects
        └── Workspaces
              ├── Settings (detail view)
              ├── Runs
              │     └── Run Detail
              ├── Variables
              │     └── Variable Detail
              ├── Config Versions
              │     ├── Config Version Detail
              │     └── Config Viewer (file tree)
              │           └── File Content
              └── State Versions
                    ├── State Version Detail
                    └── State Viewer (JSON)
```

### Workspace tabs

Once inside a workspace, you can switch between tabs using `left` and `right` arrow keys. The available tabs are:

- **Settings** -- workspace configuration details
- **Runs** -- plan and apply history
- **Variables** -- terraform and environment variables
- **Config Versions** -- uploaded configuration snapshots
- **State Versions** -- terraform state history

You can also jump directly to a tab from the workspace list using shortcut keys (`v` for variables, `f` for config versions, `s` for state versions).

## Keyboard shortcuts

Press `?` at any time to see shortcuts relevant to your current view.

### Global navigation

These keys work in every list view:

| Key | Action |
|---|---|
| `up` / `k` | Move up |
| `down` / `j` | Move down |
| `enter` | Select / drill in |
| `esc` | Go back / clear filter |
| `g` / `shift+g` | Jump to top / bottom |

### List views

| Key | Context | Action |
|---|---|---|
| `d` | Any list | View detail for selected item |
| `left` / `right` | Workspace tabs | Switch between tabs |
| `v` | Workspace list | Jump to variables tab |
| `f` | Workspace list | Jump to config versions tab |
| `s` | Workspace list | Jump to state versions tab |

### Detail and viewer views

| Key | Action |
|---|---|
| `up` / `down` | Scroll one line |
| `shift+up` / `shift+down` | Scroll half page |
| `o` | Open viewer (from state/config version detail) |

### Tools

| Key | Action |
|---|---|
| `/` | Filter the current list |
| `r` | Refresh data |
| `c` | Copy equivalent CLI command |
| `u` | Copy URL to clipboard |
| `shift+u` | Open in browser |
| `i` | Toggle instance info panel |
| `l` | Toggle API inspector |
| `?` | Toggle keyboard shortcuts |
| `q` | Quit |

## Filtering

Press `/` in any list view to activate the filter bar. Type your search term and press `enter` to apply. The list filters in real time as you type. Press `esc` to clear the filter.

Filtering works on all list views: organizations, projects, workspaces, runs, variables, config versions, and state versions.

## Scrolling in viewers

The state version JSON viewer and config version file viewer support extended scroll controls:

- `up` / `down` -- one line at a time
- `shift+up` / `shift+down` -- half a page at a time
- `r` -- re-fetch the content from the API
