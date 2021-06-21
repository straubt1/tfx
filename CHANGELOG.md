# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added

### Changed

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

