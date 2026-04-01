---
title: CLI Overview
description: Using TFx CLI commands for scripting and automation.
---

TFx CLI commands let you interact with HCP Terraform and Terraform Enterprise from scripts, CI/CD pipelines, and the terminal. Each command follows the pattern:

```sh
tfx <resource> <action> [flags]
```

For example:

```sh
tfx workspace list
tfx workspace show -w my-workspace
tfx variable list -w my-workspace
```

## Environment variables and flags

Profiles are the recommended approach, but you can also configure TFx with environment variables or flags directly. These take precedence over profile values.

| Flag | Environment Variable | Default |
|---|---|---|
| `--hostname` | `TFE_HOSTNAME` | `app.terraform.io` |
| `--default-organization` | `TFE_ORGANIZATION` | _(none)_ |
| `--token` | `TFE_TOKEN` | _(none)_ |

This is useful in CI/CD pipelines where you set credentials via environment variables rather than a config file.

## Output

CLI commands output a formatted table by default. Add `--json` (or `-j`) to get machine-readable JSON for scripting and CI pipelines.

```
$ tfx variable list -w tfx-test
╭──────────────────────┬───────────┬──────────────┬───────────┬───────┬───────────┬──────────────────────╮
│ ID                   │ KEY       │ VALUE        │ SENSITIVE │ HCL   │ CATEGORY  │ DESCRIPTION          │
├──────────────────────┼───────────┼──────────────┼───────────┼───────┼───────────┼──────────────────────┤
│ var-7XYNuuo4tMjXeXG4 │ variable7 │ {"a":"1"...} │ false     │ true  │ terraform │ I am a map in a file │
│ var-MJaLJ7czxKuU48eu │ variable3 │ It is friday │ false     │ false │ env       │ I am environmental   │
╰──────────────────────┴───────────┴──────────────┴───────────┴───────┴───────────┴──────────────────────╯
```

```sh
$ tfx variable list -w tfx-test --json | jq '.[].Key'
"variable7"
"variable3"
```
