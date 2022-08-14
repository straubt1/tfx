## Welcome to TFx Docs

Here you can find important information about TFx and how to use it.

### Commands

* [`tfx workspace`](commands/workspace.md) - Commands to work with Workspaces
* [`tfx registry`](commands/registry.md) - Commands to manage the Private Registry
* [`tfx gpg`](commands/gpg.md) - Commands to manage GPG Keys (for use with the Private Registry)
* [`tfx release`](commands/release.md) - Commands to view and download Releases

### Configuration File

Each command has the ability to pass in parameters via flags, however there are several that are required for every command.

**Example:**

```
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
tfeHostname     = "tfe.rocks" (omit to default to Terraform Cloud)
tfeOrganization = "my-awesome-org"
tfeToken        = "<Generated from Terraform Enterprise or Terraform Cloud>"
```

Common flags can also be set via environment values by using a key that is capitalized version of the flag.
This only works for the following:

- TFE_HOSTNAME
- TFE_ORGANIZATION
- TFE_TOKEN


### Disclaimer

TFx is an open source project built for use with Terraform Cloud and Terraform Enterprise under the [MIT License (MIT)](https://github.com/straubt1/tfx/blob/main/LICENSE).

!!! note ""
    While this tool is not officially supported by HashiCorp, it's current primary contributors are current or former HashiCorp employees.

### References
