---
title: Why TFx?
---

TFx exists because working with HCP Terraform and Terraform Enterprise shouldn't require writing custom API scripts for common tasks.

## The problem

The HCP Terraform and Terraform Enterprise APIs are powerful, but using them directly means writing and maintaining curl commands, Python scripts, or custom libraries for tasks that should be simple. Whether you're managing workspaces, inspecting state, reviewing runs, or automating CI/CD pipelines, the raw API adds friction.

Common challenges:

- **The web UI doesn't surface everything** -- the HCP Terraform and Terraform Enterprise UIs don't expose all API capabilities, making some operations impossible without direct API calls.
- **No single tool covers everything** -- Terraform CLI handles runs but not workspace management, variable bulk operations, registry publishing, or administrative tasks.
- **CI/CD integration is fragile** -- building platform-specific plugins for every CI system isn't feasible, and ignores the need to run the same commands locally.

## What TFx provides

TFx gives you two ways to work:

**Interactive TUI** -- run `tfx` to browse your organizations, projects, workspaces, runs, variables, state, and configuration versions in a keyboard-driven terminal interface. Ideal for exploring, investigating issues, and day-to-day operations.

**Scriptable CLI** -- use commands like `tfx workspace list` and `tfx variable list` for automation, CI/CD pipelines, and scripting. JSON output with `--json` for machine consumption.

Both share the same configuration, API client, and authentication -- use whichever fits your workflow.

## Who is TFx for?

- **Platform engineers** managing workspaces, variables, and registry modules across multiple organizations
- **DevOps teams** building CI/CD pipelines that interact with HCP Terraform or Terraform Enterprise
- **Administrators** performing bulk operations, auditing state, or managing Terraform Enterprise instances
- **Anyone** who wants to quickly browse and inspect their Terraform infrastructure without clicking through a web UI
