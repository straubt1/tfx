# Welcome to TFx Docs

_tfx_ is a standalone CLI for HCP Terraform and Terraform Enterprise.

The initial focus of _tfx_ was to execute the API-Driven workflow for a Workspace but has grown to manage multiple aspects of the platform.

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/straubt1/tfx/releases) for the latest version.

**MacOs and Linux Installation**
```sh
brew install straubt1/tap/tfx
```

**Windows Installation**

Download from the [Release](https://github.com/straubt1/tfx/releases) page.

**Go Installation**

From Go version 1.19+, the following is supported. `@latest` can be `@$VERSION`
```sh
go install github.com/straubt1/tfx@latest
```

**Container**

`tfx` is also packaged and published as a OCI Container Image in the Github Container Registry:
```sh
docker pull ghcr.io/straubt1/tfx:latest 
```

<!-- ### Commands

* [`tfx workspace`](commands/workspace.md) - Commands to work with Workspaces
* [`tfx registry`](commands/registry.md) - Commands to manage the Private Registry
* [`tfx gpg`](commands/gpg.md) - Commands to manage GPG Keys (for use with the Private Registry)
* [`tfx release`](commands/release.md) - Commands to view and download Releases -->

## Configuration File

Each command has the ability to pass in parameters via flags, however there are several that are required for every command.

**Example:**

```bash
  --tfeHostname string       The hostname of TFE without the schema. Can also be set with the environment variable TFE_HOSTNAME. (default "app.terraform.io")
  --tfeOrganization string   The name of the TFx Organization. Can also be set with the environment variable TFE_ORGANIZATION.
  --tfeToken string          The API token used to authenticate to TFx. Can also be set with the environment variable TFE_TOKEN.
```

Rather than passing these flags on each call, it is recommended to put values such as hostname, token, etc... into a configuration file.

For convenience, creating a file with the name `.tfx.hcl` and placing it in one of the following locations will auto load these values:

- Local path of where `tfx` commands are run from, `./.tfx.hcl`
- Home directory, `~/.tfx.hcl`
- Any directory you like, must specify with this path, `--config /private/somefolder/.tfx.hcl`

**Example `./.tfx.hcl`:**
```hcl
tfeHostname     = "tfe.rocks" (omit to default to HCP Terraform)
tfeOrganization = "my-awesome-org"
tfeToken        = "<Generated from Terraform Enterprise or HCP Terraform>"
```

Common flags can also be set via environment values by using a key that is capitalized version of the flag.
This only works for the following:

- TFE_HOSTNAME
- TFE_ORGANIZATION
- TFE_TOKEN

## Output Types

Most commands support a `--json` (or `-j` for short) flag that will return a proper JSON response.

> Note: Not all commands today support this flag and will ignore it.

**Default Output:**

```
$ tfx variable list -w tfx-test               
Using config file: /Users/tstraub/.tfx.hcl
List Variables for Workspace: tfx-test
╭──────────────────────┬───────────┬────────────────┬───────────┬───────┬───────────┬──────────────────────╮
│ ID                   │ KEY       │ VALUE          │ SENSITIVE │ HCL   │ CATEGORY  │ DESCRIPTION          │
├──────────────────────┼───────────┼────────────────┼───────────┼───────┼───────────┼──────────────────────┤
│ var-7XYNuuo4tMjXeXG4 │ variable7 │ {              │ false     │ true  │ terraform │ I am a map in a file │
│                      │           │   "a" = "1",   │           │       │           │                      │
│                      │           │   "b" = "zoo", │           │       │           │                      │
│                      │           │   "c" = "42"   │           │       │           │                      │
│                      │           │ }              │           │       │           │                      │
│ var-MJaLJ7czxKuU48eu │ variable3 │ It is friday   │ false     │ false │ env       │ I am environmental   │
╰──────────────────────┴───────────┴────────────────┴───────────┴───────┴───────────┴──────────────────────╯
```

**JSON Output:**

```sh
$ tfx variable list -w tfx-test --json | jq .
[
  {
    "Category": "terraform",
    "Description": "I am a map in a file",
    "HCL": true,
    "Id": "var-7XYNuuo4tMjXeXG4",
    "Key": "variable7",
    "Sensitive": false,
    "Value": "{\n  \"a\" = \"1\",\n  \"b\" = \"zoo\",\n  \"c\" = \"42\"\n}"
  },
  {
    "Category": "env",
    "Description": "I am environmental",
    "HCL": false,
    "Id": "var-MJaLJ7czxKuU48eu",
    "Key": "variable3",
    "Sensitive": false,
    "Value": "It is friday"
  }
]
```

## Disclaimer

TFx is an open source project built for use with HCP Terraform and Terraform Enterprise under the [MIT License (MIT)](https://github.com/straubt1/tfx/blob/main/LICENSE).

!!! note ""
    While this tool is not officially supported by HashiCorp, it's current primary contributors are current or former HashiCorp employees.

## References

https://github.com/hashicorp/go-tfe
https://github.com/spf13/cobra
