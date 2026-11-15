---
title: Workspace Health Assessments
---

Commands to enable or disable [Health Assessments](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/health) (drift detection and continuous validation) on Workspaces.

Health assessments are a workspace setting that can otherwise only be changed one workspace at a time in the UI. These commands make it possible to toggle the setting across many workspaces at once, for example after a bulk import, or to temporarily stop assessments from consuming agent capacity.

:::note
All commands below can be used with a `ws` alias, and `assessments` can be replaced with `health` or `drift-detection`.
:::

:::caution
Health assessments require an HCP Terraform plan or Terraform Enterprise version that supports them. If the setting is not available, the API error is reported per workspace.
:::

Workspaces that are already in the requested state are not updated and are reported as `Already enabled` / `Already disabled`.

## `tfx workspace assessments enable`

Enable health assessments on one or more workspaces by name. Multiple names can be comma separated or passed by repeating `--name`.

```sh
$ tfx workspace assessments enable --name tfx-test-workspace-01,tfx-test-workspace-02
Enabling health assessments on 2 workspace(s) in organization 'firefly'
╭───────────────────────┬─────────────────╮
│ WORKSPACE             │ STATUS          │
├───────────────────────┼─────────────────┤
│ tfx-test-workspace-01 │ Enabled         │
│ tfx-test-workspace-02 │ Already enabled │
╰───────────────────────┴─────────────────╯
```

## `tfx workspace assessments enable all`

Enable health assessments on all workspaces in an organization (sequentially).

The same filters as `tfx workspace list` can be used to select a subset of workspaces: `--search`, `--wildcard-name`, `--project-id`, `--tags`, and `--exclude-tags`.

This command will ignore individual errors and attempt to execute on all selected Workspaces.

```sh
$ tfx workspace assessments enable all --wildcard-name "tfx-test-*"
Enabling health assessments on workspaces in organization 'firefly'
Active filters:
  - wildcard-name: tfx-test-*
╭───────────────────────┬─────────────────╮
│ WORKSPACE             │ STATUS          │
├───────────────────────┼─────────────────┤
│ tfx-test-workspace-01 │ Already enabled │
│ tfx-test-workspace-02 │ Already enabled │
│ tfx-test-workspace-03 │ Enabled         │
╰───────────────────────┴─────────────────╯
```

## `tfx workspace assessments disable`

Disable health assessments on one or more workspaces by name.

```sh
$ tfx workspace assessments disable --name tfx-test-workspace-01
Disabling health assessments on 1 workspace(s) in organization 'firefly'
╭───────────────────────┬──────────╮
│ WORKSPACE             │ STATUS   │
├───────────────────────┼──────────┤
│ tfx-test-workspace-01 │ Disabled │
╰───────────────────────┴──────────╯
```

## `tfx workspace assessments disable all`

Disable health assessments on all workspaces in an organization (sequentially), optionally filtered with the same flags as `enable all`.

```sh
$ tfx workspace assessments disable all --tags agent-pool:onprem
Disabling health assessments on workspaces in organization 'firefly'
Active filters:
  - tags: agent-pool:onprem
╭───────────────────────┬──────────╮
│ WORKSPACE             │ STATUS   │
├───────────────────────┼──────────┤
│ tfx-test-workspace-04 │ Disabled │
│ tfx-test-workspace-05 │ Disabled │
╰───────────────────────┴──────────╯
```

The current setting for a single workspace is shown as `Health Assessments` in `tfx workspace show`.
