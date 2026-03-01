# Roadmap

This document tracks planned features, quality improvements, and internal work for tfx.
Priority may change — please [open an issue](https://github.com/straubt1/tfx/issues) to voice your priorities.

## Features

New commands and capabilities for users.

- [ ] Variable Sets support
- [ ] `tfx team` command group
- [ ] Cloud Agents support
- [ ] Sentinel Publishing support
- [ ] Self-signed certificate support for TFE instances
- [ ] Download workspace plan as JSON
- [ ] Plan export command
- [ ] Update `tfx organization` and `tfx project show` to include agent pool settings
- [ ] Look at brew cask and Apple developer license

### Under Consideration

- [ ] Embedded JSON filtering (similar to JMESPath in the Azure CLI)
- [ ] TUI (terminal UI) for interactive exploration
- [ ] Diff across like entities (compare two workspaces, projects, etc.)
- [ ] `-full` flag to reveal hidden global flags in help output

## Quality & Reliability

Infrastructure that improves correctness and maintainability.

- [ ] Integration testing (HCP Terraform and Terraform Enterprise)
- [ ] Automated releases with updates to the [Homebrew tap](https://github.com/straubt1/homebrew-tap)
- [ ] Automated version management (currently hardcoded in `version/version.go`)
- [ ] Consistent Go version across `go.mod`, CI workflows, and README badge
- [ ] Clean up and rename GitHub Actions workflows

## Internal / Developer

Code quality and architectural improvements.

- [ ] Add usage examples to all commands (see `tfx project` as reference)
- [ ] Spinner package (extract reusable spinner from inline usage)
- [ ] Consistent timezone display for all timestamps
- [ ] Workspace lock improvements: `lock list`, `lock show`, consider combining lock/unlock-all
- [ ] Workspace variable: move get-by-key + update/delete flow into data layer
