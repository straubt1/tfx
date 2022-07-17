<!-- <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">
 -->
# TFx CLI

_tfx_ is a standalone CLI for Terraform Cloud and Terraform Enterprise.

The initial focus of _tfx_ was to execute the API-Driven workflow for a Workspace but has grown to manage multiple aspects of the platform.

> Note: This CLI is still under active development, subject to change, and is not officially supported by HashiCorp.

[![main](https://github.com/straubt1/tfx/actions/workflows/main.yml/badge.svg)](https://github.com/straubt1/tfx/actions/workflows/main.yml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)

## Documentation

Questions about a Command? Check out the [docs](docs/docs.md)

## Why does this CLI exist?

As a consumer of Terraform Cloud or Terraform Enterprise I want to leverage the full capabilities without having to write curl/python/(insert other) libraries to call the API.

Often times these tasks are part of a delivery pipeline, but could also be administrative tasks that are done from a local machine.
The goal of this tool is to allow users to interact with the platform easily without having to create a lot of hard to maintain code.

**Common API-Driven Workflow Challenges:**

The initial use case for _tfx_ was to bridge the gap from the [CLI-Workflow](https://www.terraform.io/cloud-docs/run/cli) and the [API-Driven Workflow](https://www.terraform.io/cloud-docs/run/api).

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

## Contributing

Thank you for your interest in contributing!

_Contributing guide coming soon_

## References

https://github.com/hashicorp/go-tfe

https://github.com/spf13/cobra#installing

https://mholt.github.io/json-to-go/
