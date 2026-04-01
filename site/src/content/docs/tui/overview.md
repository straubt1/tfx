---
title: TUI Overview
description: An interactive terminal interface for browsing and managing HCP Terraform and Terraform Enterprise resources.
---

The TFx TUI is an interactive terminal interface for browsing your HCP Terraform or Terraform Enterprise resources. Instead of running individual CLI commands, you can navigate organizations, projects, workspaces, runs, variables, state, and configuration versions in a single session.

## Launching the TUI

Run `tfx` with no arguments:

```sh
tfx
```

The TUI uses the same configuration as the CLI. If you haven't configured TFx yet, see [Getting Started](/gettingstarted/).

To authenticate for the first time or manage profiles:

```sh
tfx login
```

<!-- ![TFx TUI demo](/media/tui-demo.gif) -->

## What you'll see

When the TUI starts, you land on the **Organizations** list showing all organizations accessible to your API token. From there, drill into any resource:

```
Organizations  >  Projects  >  Workspaces  >  Runs / Variables / State / Config
```

A breadcrumb bar at the top always shows where you are. Press `esc` to go back, `enter` to drill in.

## Quick reference

| Key | Action |
|---|---|
| `enter` | Select / drill in |
| `esc` | Go back |
| `/` | Filter the current list |
| `?` | Show keyboard shortcuts |
| `q` | Quit |

See [Navigation](/tui/navigation/) for the full keyboard reference and [Features](/tui/features/) for details on the API inspector, state viewer, and more.
