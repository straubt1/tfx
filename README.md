<!-- <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">
 -->
# TFx CLI

_tfx_ is a standalone CLI for Terraform Cloud and Terraform Enterprise.

The initial focus of _tfx_ is to execute the API-Driven workflow for a Workspace, but will expand to other common workflows that, in the past, have required API wrappers.

> Note: This CLI is still under active development, subject to change, and not officially supported by HashiCorp.

[![main](https://github.com/straubt1/tfx/actions/workflows/main.yml/badge.svg)](https://github.com/straubt1/tfx/actions/workflows/main.yml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)

## Why does this CLI exist?

As a consumer of TFC/TFE I want to leverage the full capabilities without having to write curl/python/(insert other) libraries to call the API.
Often times these tasks are part of my pipeline, but could also be administrative tasks that are done from a local machine.

**Common API-Driven Workflow Challenges:**

- The CLI-Driven workflow presents several challenges when creating more advanced pipelines for a Workspace run, specifically the inability to insert a gate check between a plan and apply, (in other words you must run a `terraform apply -auto-approve`).
- The CLI driven workflow requires a `terraform init` that forces a download of providers before a plan can be called remotely, these providers are never actually used on the local host and can be difficult to source airgapped environments.
- Implementing an API-Driven workflow requires several API calls to perform a plan/apply.
- It is unlikely that the full range of features will be built into [Terraform](https://github.com/hashicorp/terraform).
- Developing CI/CD specific plugins for even the most common tools is not feasible, and ignores the ability to run the commands locally.

![Terminal Example Plan](assets/terminal-example-plan.gif)
## Installation

Binaries are created as part of release, check out the [Release Page](https://github.com/straubt1/tfx/releases) for the latest release.

**MacOs Installation**
```sh
version="0.0.3-dev"
curl -L -o tfx "https://github.com/straubt1/tfx/releases/download/${version}/tfx_darwin_amd64"
chmod +x tfx
```

**Linux Installation**
```sh
version="0.0.3-dev"
curl -L -o tfx "https://github.com/straubt1/tfx/releases/download/${version}/tfx_linux_amd64"
chmod +x tfx
```

**Windows Installation**
```sh
version="0.0.3-dev"
curl -L -o tfx.exe "https://github.com/straubt1/tfx/releases/download/${version}/tfx_windows_amd64"
```

**Go Installation**
From Go version 1.18, the following is supported. `@latest` can be `@$VERSION`
```sh
go install github.com/straubt1/tfx@latest
```

## Usage

Each command has the ability to pass in parameters via flags, several are required for every command.

Example:
```
  --tfeHostname string       The hostname of TFE without the schema. Can also be set with the environment variable TFE_HOSTNAME. (default "app.terraform.io")
  --tfeOrganization string   The name of the TFx Organization. Can also be set with the environment variable TFE_ORGANIZATION.
  --tfeToken string          The API token used to authenticate to TFx. Can also be set with the environment variable TFE_TOKEN.
```

Flags can also be created in a configuration file with the file name ".tfx.hcl".
Flags can also be set via environment values by using a key that is capitalized version of the flag.

For convenience this file will automatically load if it is in the hosts home directory or current working directory.

Example:
`./.tfx.hcl`
```hcl
tfeHostname     = "tfe.rocks" (omit to default to Terraform Cloud)
tfeOrganization = "my-awesome-org"
tfeToken        = "<Generated from Terraform Enterprise or Terraform Cloud>"
```

You can also specify this file via the `--config` flag.

## Workspace Run Workflow

**Create a Plan**

```sh
# Create a speculative plan that can not be applied
tfx plan -w tfx-test -s

# Create a plan that can be applied
tfx plan -w tfx-test

# Create a Configuration Version based on terraform in the current directory
tfx cv create -w tfx-test

# Create a Configuration Version based on terraform in a supplied directory
tfx cv create -w tfx-test -d ./myterraformfolder/

# Create a plan based on a configuration version
tfx plan -w tfx-test -i cv-HKE8gevVtGBXapcq
```

**Create an Apply**

```sh
tfx apply -r <run-id>
```

## Commands

### `tfx plan`

Create a plan to execute on TFx.

`tfx plan` - Create a Workspace plan based on a current directory

`tfx plan export` - Create and download a Sentinel Mock from a plan

### `tfx apply`

Create an apply to execute on TFx.

`tfx apply` - Apply a Workspace plan based on a plan

### `tfx run`

Managing Workspace Runs.

`tfx run list` - List all Runs for a supplied Workspace

`tfx run create` - Create a Run for a supplied Workspace

`tfx run show` - Show Run details for a supplied Run

### `tfx cv`

Managing Workspace Configuration Versions.

`tfx cv list` -  List all Configuration Versions for a supplied Workspace

`tfx cv create` - Create a Configuration Version for a supplied Workspace

`tfx cv show` - Show Configuration Version details for a supplied Configuration

### `tfx pmr`

Managing Private Module Registry modules.

`tfx pmr list` - List all modules in the PMR

`tfx pmr create` - Create a module in the PMR

`tfx pmr create version` - Create a version of a module in the PMR

`tfx pmr show` - Show module details for a supplied module

`tfx pmr show versions` - Show modules versions for a supplied module

`tfx pmr delete` - Delete a module from the PMR

`tfx pmr delete version` - Delete a specific module version from the PMR

`tfx pmr download` - Download a specific module version of TF code

### `tfx state`

Managing Workspace State Files (State Versions).

`tfx state list` - List all State Versions for a supplied Workspace

`tfx state show` - Show state details for a supplied State Version

`tfx state download` - Download a specific State Version

`tfx state create` - Create a new State Version with a supplied state file
- There is no way to delete State Versions
- The LAST State Version to be created is the "current" state file that will be used by the Workspace
- A Workspace must be locked to create new State Versions
- The "serial" attribute must be incremented
- The "lineage" attribute must be the same for any newly created State Version
- The API does not return a state versions lineage, you must download the file and parse to get the lineage

### `tfx tfv`

Managing Terraform Versions in a Terraform Enterprise install (TFE only).

`tfx tfv list` - List all Terraform Versions for a Terraform Enterprise install

`tfx tfv show` - Show version details for a supplied Terraform Version or Version Id

`tfx tfv delete` - Delete a version of a supplied Terraform Version or Version Id

`tfx tfv create` - Create a Terraform Version

`tfx tfv create official` - Create an official Terraform Version from releases.hashicorp.com

`tfx tfv disable` - Disable a Terraform Version(s), accepts comma separated list

`tfx tfv disable all` - Disable all Terraform Versions

`tfx tfv enable` - Enables a Terraform Version(s), accepts comma separated list

`tfx tfv enable all` - Enables all Terraform Versions

### `tfx workspace`

`tfx workspace list`
- Using the "--search" flag allows filtering by workspaces with a given string
- Using the "--run-status" flag allows filtering by workspaces with a current run with a given status (full list of available run statuses can be found [here](https://www.terraform.io/docs/cloud/api/run.html#run-states))

`tfx workspace list all`
- Using the "--search" flag allows filtering by workspaces with a given string

`tfx workspace lock`
- Lock a given workspace by name, in a given organization

`tfx workspace lock all`
- Lock all workspaces in a given organization

`tfx workspace unlock`
- Unlock a given workspace by name, in a given organization

`tfx workspace unlock all`
- Unlock all workspaces in a given organization

### `tfx variable`

Manage Workspace Variables

`tfx variable list`
- List all workspace variables

`tfx variable create`
- Create a workspace variable, optionally the value can read from a filename

`tfx variable update`
- Update an existing workspace variable, optionally the value can read from a filename

`tfx variable show`
- Show details of a workspace variable

`tfx variable delete`
- Delete a workspace variable 


### `tfx release`

Manage [Terraform Enterprise](https://www.terraform.io/enterprise) and [Replicated](https://www.replicated.com/) releases and binaries needed for airgap installations (TFE only).

`tfx release tfe list`
- List available Terraform Enterprise releases

`tfx release tfe show`
- Show details of a Terraform Enterprise release, including release notes

`tfx release tfe download`
- Download a Terraform Enterprise airgap binary

`tfx release replicated list`
- List available Replicated releases

`tfx release replicated download`
- Download a Replicated release

### `tfx registry provider`

Manage providers in the Private Registry

`tfx registry provider list`
- List Providers in the Registry

`tfx registry provider version list`
- List Versions for a Provider in the Registry

`tfx registry provider version platform list`
- List Platforms for a Provider Version in the Registry

`tfx registry provider create`
- Create a Provider in the Registry

`tfx registry provider show`
- Show details of a Provider in the Registry

`tfx registry provider delete`
- Delete a Provider in the Registry

`tfx registry provider version create`
 - Create a Version for a Provider in the Registry

`tfx registry provider version show`
- Show details a Version for a Provider in the Registry

`tfx registry provider version delete`
 - Delete a Version for a Provider in the Registry

`tfx registry provider version platform create`
- Create a Platform Version for a Provider in the Registry

`tfx registry provider version platform show`
- Show details of a Platform Version for a Provider in the Registry

`tfx registry provider version platform delete`
- Delete a Platform Version for a Provider in the Registry

## Potential Future Commands

Additional commands to implement.

- [ ] `tfx run`
  - [ ] `cancel`, cancel, discard, force cancel a run
- [ ] `tfx pmr`
  - [ ] `search` find a module https://www.terraform.io/docs/registry/api.html#search-modules
- [ ] `tfx sentinel`
  - [ ] `list`, list policy sets
  - [ ] `create`, create a policy set
  - [ ] `delete`, deletes a policy set
  - [ ] `assign`, assigns a WS(s) to the policy set

## Contributing

Thank you for your interest in contributing!

_Contributing guide coming soon_

## References

https://github.com/hashicorp/go-tfe

https://github.com/spf13/cobra#installing

https://mholt.github.io/json-to-go/
