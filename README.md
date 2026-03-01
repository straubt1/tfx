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
| `tfx workspace configurationversion` | `list`, `show`, `create`, `download` |
| `tfx workspace lock` | `lock`, `unlock`, `list` |
| `tfx workspace run` | `list`, `show`, `create`, `cancel`, `discard` |
| `tfx workspace stateversion` | `list`, `show`, `create`, `download` |
| `tfx workspace team` | `list` |
| `tfx workspace variable` | `list`, `show`, `create`, `update`, `delete` |
| `tfx plan` | `show`, `logs`, `jsonoutput`, `create` |
| `tfx registry module` | `list`, `show`, `upload`, `delete` |
| `tfx registry module-version` | `list`, `show`, `delete` |
| `tfx registry provider` | `list`, `show`, `upload`, `delete` |
| `tfx registry provider-version` | `list`, `show`, `delete` |
| `tfx registry provider-version-platform` | `list`, `show`, `delete` |
| `tfx release tfe` | `list`, `show` |
| `tfx admin gpg` | `list`, `show`, `create`, `delete` |
| `tfx admin metrics` | `workspace` |
| `tfx admin terraformversion` | `list`, `show`, `create`, `update`, `delete` |

## Roadmap

Future implementation items:

- [ ] Support Variable Sets
- [ ] Support Sentinel Publishing
- [ ] Support Cloud Agents

Priority may change, please submit an issue to voice changes you would like to see.
