# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

* `tfx workspace lock all` and `tfx workspace unlock all` now support a `--search` argument to filter affected workspaces by name search.

### Changed

* Moved and updated docs to the `site/` folder and published to tfx.rocks custom domain.
* Updated all of the `tfx registry module` commands to support JSON (`--json`) output.
* Moved Module Registry commands from the `pmr` command, the new commands are:
  * `tfx registry module list`
  * `tfx registry module create`
  * `tfx registry module show`
  * `tfx registry module delete`
  * `tfx registry module version list`
  * `tfx registry module version create`
  * `tfx registry module version delete`
  * `tfx registry module version download`
* Updated all of the `tfx registry provider` commands to support JSON (`--json`) output.
* Updated `tfx workspace` commands to support JSON (`--json`) output.
  * `tfx workspace list`
  * `tfx workspace show`
  * `tfx workspace lock`
  * `tfx workspace lock all`
  * `tfx workspace unlock`
  * `tfx workspace unlock all`
* Updated all of the `tfx workspace configuration-version` commands to support JSON (`--json`) output.
* Moved Configuration Version commands from the `cv` command, there is an alias for `cv` to make typing commands easier, the new commands are:
  * `tfx workspace configuration-version list` or `tfx workspace cv list`
  * `tfx workspace configuration-version create` or `tfx workspace cv create`
  * `tfx workspace configuration-version show` or `tfx workspace cv show`
  * `tfx workspace configuration-version download` or `tfx workspace cv download`
* Updated all of the `tfx workspace state-version` commands to support JSON (`--json`) output.
* Moved State Version commands from the `state` command, there is an alias for `sv` to make typing commands easier, the new commands are:
  * `tfx workspace state-version list` or `tfx workspace sv list`
  * `tfx workspace state-version create` or `tfx workspace sv create`
  * `tfx workspace state-version show` or `tfx workspace sv show`
  * `tfx workspace state-version download` or `tfx workspace sv download`
* Moved Module Registry commands from the `tfv` command, the new commands are:
  * `tfx admin terraform-version list`
  * `tfx admin terraform-version create`
  * `tfx admin terraform-version create official`
  * `tfx admin terraform-version show`
  * `tfx admin terraform-version delete`
  * `tfx admin terraform-version enable`
  * `tfx admin terraform-version enable all`
  * `tfx admin terraform-version disable`
  * `tfx admin terraform-version disable all`
* Updated all of the `tfx admin terraform-version` commands to support JSON (`--json`) output.

### Removed


## [0.0.4-dev] - 2022.07.17

### Added

New Commands:

* `tfx cv download` - Download a Configuration Version and unpack onto disk.
* `tfx gpg list` - List GPG Keys of an Organization
* `tfx gpg create` - Create a GPG Key for an Organization
* `tfx gpg show` - Show details of a GPG Key for an Organization
* `tfx gpg delete` - Delete GPG Key for an Organization
* `tfx release tfe list` - List available Terraform Enterprise releases
* `tfx release tfe show` - Show details of a Terraform Enterprise release, including release notes
* `tfx release tfe download` - Download a Terraform Enterprise airgap binary
* `tfx release replicated list` - List available Replicated releases
* `tfx release replicated download` - Download a Replicated release
* `tfx registry provider list` - List Providers in the Registry
* `tfx registry provider version list` - List Versions for a Provider in the Registry
* `tfx registry provider version platform list` - List Platforms for a Provider Version in the Registry
* `tfx registry provider create` - Create a Provider in the Registry
* `tfx registry provider show` - Show details of a Provider in the Registry
* `tfx registry provider delete` - Delete a Provider in the Registry
* `tfx registry provider version create` - Create a Version for a Provider in the Registry
* `tfx registry provider version show` - Show details a Version for a Provider in the Registry
* `tfx registry provider version delete` - Delete a Version for a Provider in the Registry
* `tfx registry provider version platform create` - Create a Platform Version for a Provider in the Registry
* `tfx registry provider version platform show` - Show details of a Platform Version for a Provider in the Registry
* `tfx registry provider version platform delete` - Delete a Platform Version for a Provider in the Registry
* `tfx variable list` - List all workspace variables
* `tfx variable create` - Create a workspace variable, optionally the value can read from a filename
* `tfx variable update` - Update an existing workspace variable, optionally the value can read from a filename
* `tfx variable show` - Show details of a workspace variable
* `tfx variable delete` - Delete a workspace variable

### Changed

* Upgraded to Go 1.18.4
* PMR Module uploading - Removed helper code (shim) in favor of the now available the go-tfe functions
* Lots of refactoring within some commands
* Added optional `--json` flag framework to allow output to be in JSON for non-interactive use
* Started the migrating to "github.com/fatih/color" from "github.com/logrusorgru/aurora" for text colors

### Removed

## [0.0.3-dev] - 2021.06.22

### Added

* `tfx workspace lock` - Lock a given workspace by name, in a given organization
* `tfx workspace lock all` - Lock all workspaces in a given organization
* `tfx workspace unlock` - Unlock a given workspace by name, in a given organization
* `tfx workspace unlock all` - Unlock all workspaces in a given organization

### Changed

* `tfx workspace` commands now sort WS by name

### Removed

## [0.0.2-dev] - 2021.06.20

### Added

* hostname, organization and token can now be set with the respective environment values to align with [TFE Provider](https://registry.terraform.io/providers/hashicorp/tfe/latest/docs). ([#7](https://github.com/straubt1/tfx/issues/7))
  * TFE_HOSTNAME
  * TFE_ORGANIZATION
  * TFE_TOKEN
* Added "message" flag to `tfx run` and `tfx plan` commands. ([#8](https://github.com/straubt1/tfx/issues/8))
* `tfx workspace` commands
  * `list` - List all workspaces in an Organization (optional workspace name search string) 
  * `list all` - List all workspaces in All Organizations the API token has access to (optional workspace name search string) 
  * `show` - Show details of a workspace
* `tfx metrics` command to pull details about TFx (this command is hidden)
  * Organization Count
  * Workspace Count
  * Run Count
  * Policy Check Count
  * Policies Pass/Fail Count
* `tfx metrics workspace` command to get run metrics for all workspaces in a single organization (this command is hidden)
  * Can filter on start date
  * Output:
    * Workspace Name
    * Total Runs
    * Errored Runs
    * Discarded Runs
    * Cancelled Runs

### Changed

* Cleaned up docs

### Removed

## [0.0.1-dev] - 2021.05.23

### Added

* `tfx plan export` command to download sentinel mock data
* `tfx state` commands
  * list
  * show
  * download
  * create
* `tfx tfv` commands
  * list
  * show
  * create
  * create official
  * delete
  * disable
  * disable all
  * enable
  * enable all

### Changed

* `tfx plan export` added flag to supply a directory
* Added "Built By" output on version to help originate a build

### Removed

## [0.0.0-dev] - 2021.05.16

### Added

### Changed

### Removed


[Unreleased]: https://github.com/straubt1/tfx/compare/v1.0.0...HEAD
[0.0.1]: https://github.com/ostraubt1/tfx/compare/v0.0.0...v0.0.1 
[0.0.0]: https://github.com/straubt1/tfx/releases/tag/v0.0.1

