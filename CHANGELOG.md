# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.1] - Unreleased

**Changed**

* Removed `-p` shorthand from `--project-id` (workspace list) and `--provider` (registry module commands) to avoid conflict with the root `--profile` (`-p`) flag

## [v0.3.0] - 2026-03-31

**Added**

* `tfx` now launches the interactive TUI by default — no subcommand needed
* Interactive TUI: full-screen terminal browser for HCP Terraform and TFE
  * Navigate organizations → projects → workspaces → runs / variables / config versions / state versions
  * Detail views for all resource types with scrollable content
  * State version JSON viewer with syntax highlighting and line numbers
  * Config version file browser: download and browse uploaded files with a content viewer
  * Workspace tabs (Settings, Runs, Variables, Config Versions, State Versions) with `←`/`→` switching
  * API Inspector panel (`l`): live capture of every API call with method, path, status, duration; drill into full response bodies; copy response or curl command
  * Instance Info modal (`i`): hostname, API version, TFE version, health check status
  * Profile bar showing active profile, username, email, and token expiry
  * Context-aware help overlay (`?`): shows only shortcuts relevant to the current view
  * List filtering (`/`) on all views
  * Copy CLI command (`c`), copy URL (`u`), open in browser (`shift+u`)
* `tfx login` — interactive authentication wizard with profile management, token validation, and org selection
* Profile-based config format in `.tfx.hcl` with named profile blocks and multi-profile support
* `tfx workspace run policy` — policy check and evaluation details for a run
* `--logs` flag for `tfx workspace run policy` — raw policy output (Sentinel/OPA)

**Fixed**

* State version JSON viewer: consistent background colors on indentation and whitespace
* API Inspector: background colors now match panel focus state
* Login token validation: stricter format check detects modified/truncated tokens
* CLI hint bar shows correct command syntax for all detail views
* Multi-line variable values collapsed to `↵` in list view, expanded in detail view

**Changed**

* `tfx tui` subcommand removed — `tfx` launches the TUI directly; all CLI commands unchanged
* Config keys renamed: `tfeHostname` → `hostname`, `tfeOrganization` → `defaultOrganization`, `tfeToken` → `token` (old keys still supported)
* Keyboard shortcut labels use `shift+u`/`shift+c`/`shift+g` instead of uppercase letters
* Documentation site redesigned: new Interactive TUI section, CLI Commands overview, updated Getting Started with `tfx login` flow, "Why TFx?" rewrite, accent-tinted inline code styling

## [v0.2.1] - 2026-03-01

**Added**

* `tfx variable-set list` — list variable sets in an organization, with optional `--search` filter (alias: `varset`)
* `tfx variable-set show` — show details of a variable set by ID, including assigned workspaces, projects, and variables
* `tfx variable-set create` — create a variable set with optional `--description`, `--global`, and `--priority` flags
* `tfx variable-set delete` — delete a variable set by ID

**Fixed**

* `tfx ws cv download` no longer fails with "no such file or directory" when no `--directory` flag is provided; the target subdirectory is now created automatically
* Updated documentation site
* Created a Test Plan on the website for those who want to run through the commands and test out the functionality

## [v0.2.0] - 2026-03-01

**Added**

* `tfx organization list` — list all organizations the token has access to
* `tfx organization show` — show details of an organization
* `tfx workspace plan show` — show details of a plan by ID
* `tfx workspace plan logs` — stream logs for a plan by ID
* `tfx workspace plan jsonoutput` — show JSON execution output for a plan by ID
* `tfx workspace plan create` — create a plan for a workspace

**Changed**

* `tfx workspace run apply` now supports additional run statuses (`planned_and_saved`, `cost_estimated`, `policy_checked`)
* Plan commands moved under `tfx workspace plan` for consistency (previously `tfx plan`)
* Major internal refactor: layered architecture with dedicated `client/`, `data/`, `output/`, and `cmd/views/` packages

**Removed**

* `tfx plan` and `tfx apply` commands - These may be added back in the future.
* `tfx release replicated` commands — TFE no longer runs on Replicated; use a previous version of `tfx` if needed
* `tfx release tfe download` — TFE now uses Docker images; use `docker pull` instead

## [v0.1.5] - 2025-09-06

**Added**

* Added `tfx project list`
* Added `tfx project show --id`
* Added Resource Count with `tfx workspace list` and `tfx workspace show` commands

**Changed**

* Removed HCL Configuration Option
* Better Config File error handling (and update to maintain .hcl configuration files)
* Updated references to "HCP Terraform"

## [v0.1.4] - 2024-10-29

**Added**

* Added `tfx workspace run cancel` to cancel the latest run.

**Changed**

* Updated all go packages to latest
* Updated documentation example errors
* Fixed BuiltBy build tag for releases

## [v0.1.3] - 2024-09-20

**Added**

**Changed**

* Updated go to v1.23
* Updated all go packages to latest
* Fixed Workspace Variable Create when the value passed in would error in the `isFile()` function.
* Bumped Docker to alpine:3.20.3

**Removed**

* `tfx workspace list --repository` flag has been removed.

## [v0.1.2] - 2023-10-13

**Added**

* `tfx workspace run discard` - Discard a run on a given Workspace.

**Changed**

* Updated go to v1.21
* Updated all go packages to latest
* `tfx admin gpg create` and `tfx admin gpg show` commands will now output "KeyID" instead of "ID" to match the API response.

## [v0.1.1] - 2022-09-11

**Added**

* `tfx workspace team list` - List team access information for a given Workspace.

**Changed**

* Updated `tfx workspace list` to add a flag `--repository` that will filter based on workspaces that are connected to the specified repository identifier.

* `tfx workspace show` - Now includes Team Access and Statefile Sharing output for the Workspace.

**Removed**

## [v0.1.0] - 2022-08-21

First official release!

**Added**

* `tfx workspace lock all` and `tfx workspace unlock all` now support a `--search` argument to filter affected workspaces by name search.

**Changed**

* Moved and updated docs to the `site/` folder and published to a custom domain [tfx.rocks](tfx.rocks).
* Created a style guide and updated Commands, more information can be found at [https://tfx.rocks/about/style_guide/]()
  * Some of these changes are **BREAKING** changes
  * Mainly moving away from Command Flags that used camel case (example: "workspaceName" to "workspace-name")
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

**Removed**


## [0.0.4-dev] - 2022-07-17

**Added**

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

**Changed**

* Upgraded to Go 1.18.4
* PMR Module uploading - Removed helper code (shim) in favor of the now available the go-tfe functions
* Lots of refactoring within some commands
* Added optional `--json` flag framework to allow output to be in JSON for non-interactive use
* Started the migrating to "github.com/fatih/color" from "github.com/logrusorgru/aurora" for text colors

**Removed**

## [0.0.3-dev] - 2021-06-22

**Added**

* `tfx workspace lock` - Lock a given workspace by name, in a given organization
* `tfx workspace lock all` - Lock all workspaces in a given organization
* `tfx workspace unlock` - Unlock a given workspace by name, in a given organization
* `tfx workspace unlock all` - Unlock all workspaces in a given organization

**Changed**

* `tfx workspace` commands now sort WS by name

**Removed**

## [0.0.2-dev] - 2021-06-20

**Added**

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

**Changed**

* Cleaned up docs

**Removed**

## [0.0.1-dev] - 2021-05-23

**Added**

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

**Changed**

* `tfx plan export` added flag to supply a directory
* Added "Built By" output on version to help originate a build

**Removed**

## [0.0.0-dev] - 2021-05-16

**Added**

**Changed**

**Removed**

[Unreleased]: https://github.com/straubt1/tfx/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/straubt1/tfx/releases/tag/v0.1.0
