---
title: Features
description: API inspector, state viewer, config viewer, and other TUI features.
---

## Profile bar

The top of the TUI displays your active profile information:

- **profile** -- the named profile from your config file (e.g., `default`)
- **username** -- your authenticated user
- **email** -- account email
- **expires** -- token expiration date (or `never` for non-expiring tokens)

This data is fetched from the API on startup using your configured token.

## API Inspector

Toggle the API inspector with `l`. This opens a split panel on the right side showing every API call the TUI makes in real time.

The inspector has two modes:

**List mode** -- shows all API requests with method, path, status code, and duration. Use `up`/`down` to browse and `/` to filter.

**Detail mode** -- press `enter` on any request to see the full response body with syntax-highlighted JSON. From here:

- `c` -- copy the response body to your clipboard
- `shift+c` -- copy an equivalent `curl` command to your clipboard
- `esc` -- return to the list

Press `tab` to switch focus between the main content and the inspector panel.

## State Version Viewer

From the state versions tab, press `enter` on any state version to open the JSON viewer. The viewer displays the full Terraform state with syntax highlighting:

- **Keys** in blue
- **Strings** in green
- **Numbers** in purple
- **Booleans and null** in amber
- **Punctuation** in gray

Line numbers are shown on the left. Scroll with `up`/`down` or `shift+up`/`shift+down` for half-page jumps. Press `r` to re-fetch the state from the API.

State files are cached locally at `~/.tfx/cache/state/` so subsequent views load instantly.

## Config Version Viewer

From the config versions tab, press `enter` to open the file tree browser. This shows all files uploaded in the configuration version.

Navigate the tree with `up`/`down` and press `enter` to view a file's contents. The file content viewer supports the same scroll controls as the state viewer.

Press `p` to copy the local cache path to your clipboard. Config files are cached at `~/.tfx/cache/cv/`.

## Instance Info

Press `i` to toggle the instance info panel. This shows details about the connected TFE/HCP Terraform instance including version and health status. Press `r` while the panel is open to refresh.

## Copy CLI Command

Press `c` in any list view to copy the equivalent `tfx` CLI command to your clipboard. For example, while viewing a workspace list, pressing `c` copies a command like:

```
tfx workspace list --default-organization my-org
```

This makes it easy to transition from interactive browsing to scripted automation.

## URL Actions

- `u` -- copy the URL of the selected resource to your clipboard (works for workspaces, runs, organizations, projects, and more)
- `shift+u` -- open the selected resource directly in your default browser
