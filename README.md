# TFx CLI

_tfx_ is a standalone CLI for Terraform Cloud and Terraform Enterprise.

The initial focus of _tfx_ is to execute the API-Driven workflow for a workspace, but will expand to other common workflows that, in the past, have required API wrappers.

## Why does this CLI exist?

As a consumer of TFC/TFE I want to leverage the full capabilities without having to write curl/python/(insert other) libraries to call the API.
Often times these tasks are part of my pipeline, but could also be administrative tasks that are done from a local machine.

**Common API-Driven Workflow Challenges:**

- The CLI-Driven workflow presents several gaps in creating more advanced pipelines a Workspace run, specifically the inability to insert a gate check between a plan and apply, (in other words you must run a `terraform apply -auto-approve`).
- The CLI driven workflow requires a `terraform init` that forces a download of providers before a plan can be called remotely, these providers are never actually used on the local host and can be difficult to source in airgap environments.
- Implementing an API-Driven workflow requires several API calls to perform a plan/apply.
- It is unlikely that the full range of features will be built into [Terraform](https://github.com/hashicorp/terraform).
- Developing CI/CD specific plugins for even the most common tools is not feasible, and ignores the ability to run the commands locally.

## Installation

Binaries are created as part of release, check out the [Release Page](https://github.com/straubt1/tfx/releases) for the latest release.

**MacOs Installation**
```sh
version="0.0.0"
curl -L -o tfx "https://github.com/straubt1/tfx/releases/download/${version}/tfx_darwin_amd64"
chmod +x tfx
```

**Linux Installation**
```sh
version="0.0.0"
curl -L -o tfx "https://github.com/straubt1/tfx/releases/download/${version}/tfx_linux_amd64"
chmod +x tfx
```

**Windows Installation**
```sh
version="0.0.0"
curl -L -o tfx.exe "https://github.com/straubt1/tfx/releases/download/${version}/tfx_windows_amd64"
```

## Setup

Each command has the ability to pass in parameters via flags.

Example:
```
  --tfeHostname string       The hostname of TFE without the schema (defaults to TFE app.terraform.io). (default "app.terraform.io")
  --tfeOrganization string   The name of the TFx Organization.
  --tfeToken string          The API token used to authenticate to TFx.
```

Flags can also be created in a configuration file with the file name ".tfx.hcl".
For convenience this file will automatically load if it is in the hosts home directory or current working directory.

Example:
`./.tfx.hcl`
```hcl
tfeHostname     = "tfe.rocks" (omit to default to TFC)
tfeOrganization = "my-awesome-org"
tfeToken        = "<Generated from TFx>"
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

## Future Commands

- [ ] `tfx run`
  - [x] `list`, List all runs for a workspace
  - [x] `create`, Create a Run
  - [x] `show`, Show a run
  - [ ] `cancel`, cancel, discard, force cancel a run
- [x] `tfx cv`
  - [x] `list`, List all configuration versions for a workspace
  - [x] `create`, Create a configuration version
  - [x] `show`, Show a configuration version
- `tfe tfversions`
  - [ ] `list`, list all Terraform versions in TFE
  - [ ] `disable`, disable a Terraform version, -a flag to disable all
  - [ ] `enable`, enable a Terraform version
  - [ ] `create`, create a new Terraform version, upsert?
  - [ ] `show`, show a version
  - [ ] `delete`, delete a version
- `tfx pmr`
  - [x] `list`, list all modules in the PMR
  - [x] `create`, create a module in the PMR
  - [x] `create version`, create a version of a module
  - [x] `show`, show a module
  - [x] `show versions`, show a modules versions
  - [x] `delete`, deletes a module from the PMR
  - [x] `delete version`, deletes a module version from the PMR
  - [x] `download` download a version of TF code
  - [ ] `search` find a module https://www.terraform.io/docs/registry/api.html#search-modules
- [ ] `tfe sentinel`
  - [ ] `list`, list policy sets
  - [ ] `create`, create a policy set
  - [ ] `delete`, deletes a policy set
  - [ ] `assign`, assigns a WS(s) to the policy set


## References

https://github.com/hashicorp/go-tfe

https://github.com/spf13/cobra#installing

https://mholt.github.io/json-to-go/

