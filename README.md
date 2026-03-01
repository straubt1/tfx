# HCP Terraform/Terraform Enterprise CLI

[![main](https://github.com/straubt1/tfx/actions/workflows/main.yml/badge.svg)](https://github.com/straubt1/tfx/actions/workflows/main.yml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.25-61CFDD.svg?style=flat-square)

_tfx_ is a standalone CLI for HCP Terraform and Terraform Enterprise.

The initial focus of _tfx_ was to execute the API-Driven workflow for a Workspace but has grown to manage multiple aspects of the platform.

> Note: This CLI is still under active development, subject to change, and is not officially supported by HashiCorp.

## Documentation

Looking for more information?
Check out our docs site [tfx.rocks](https://tfx.rocks)

## Commands

| Command Group | Subcommands |
|---|---|
| `tfx organization` | `list`, `show` |
| `tfx project` | `list`, `show` |
| `tfx workspace` | `list`, `show` |
| `tfx workspace configuration-version` | `list`, `show`, `create`, `download` |
| `tfx workspace lock` | `all` |
| `tfx workspace unlock` | `all` |
| `tfx workspace plan` | `show`, `logs`, `jsonoutput`, `create` |
| `tfx workspace run` | `list`, `show`, `create`, `cancel`, `discard` |
| `tfx workspace state-version` | `list`, `show`, `create`, `download` |
| `tfx workspace team` | `list` |
| `tfx workspace variable` | `list`, `show`, `create`, `update`, `delete` |
| `tfx registry module` | `list`, `show`, `create`, `delete` |
| `tfx registry module version` | `list`, `create`, `delete`, `download` |
| `tfx registry provider` | `list`, `show`, `create`, `delete` |
| `tfx registry provider version` | `list`, `show`, `create`, `delete` |
| `tfx registry provider version platform` | `list`, `show`, `create`, `delete` |
| `tfx release tfe` | `list`, `show` |
| `tfx admin gpg` | `list`, `show`, `create`, `delete` |
| `tfx admin metrics` | `workspace` |
| `tfx admin terraformversion` | `list`, `show`, `create`, `update`, `delete` |

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and improvements.
Priority may change — [open an issue](https://github.com/straubt1/tfx/issues) to voice your priorities.
