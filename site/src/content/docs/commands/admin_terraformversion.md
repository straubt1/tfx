---
title: Terraform Version Commands
---

Managing Terraform Versions in a Terraform Enterprise install.

:::note
These commands will only work with Terraform Enterprise
:::

:::note
All commands below can be used with a `tfv` alias.
:::

## Command overview

| Subcommand | Description |
|---|---|
| `list` | List all Terraform versions (`--search` for substring filter) |
| `show` | Show details for one version |
| `create` | Create a custom version from URL and SHA |
| `create official` | Create a version from HashiCorp releases |
| `delete` | Delete a version |
| `disable` | Disable specific versions (`--versions`) |
| `disable all` | Bulk disable; optional filter flags select a subset |
| `enable` | Enable specific versions (`--versions`) |
| `enable all` | Bulk enable; optional filter flags select a subset |

Bulk `disable all` and `enable all` accept at most one filter flag per invocation. See the sections below for flag details and examples.

## `tfx admin terraform-version list`

List all Terraform Versions for a Terraform Enterprise install.

Using the `--search` flag allows filtering by version contains a given string.

**Example**

```sh
$ tfx admin terraform-version list
Using config file: /Users/tstraub/.tfx.hcl
List Terraform Versions for TFE 
╭───────────────────────┬───────────────────────┬─────────┬──────────┬───────┬────────────╮
│ VERSION               │ ID                    │ ENABLED │ OFFICIAL │ USAGE │ DEPRECATED │
├───────────────────────┼───────────────────────┼─────────┼──────────┼───────┼────────────┤
│ 1.2.1                 │ tool-bcVckVA7aH7b98WR │ true    │ true     │     0 │ false      │
│ 1.2.0                 │ tool-7otJm3yaad55iBQY │ true    │ true     │     0 │ false      │
│ 1.2.0-rc2             │ tool-1TgFEgk2nSFfC8e4 │ true    │ true     │     0 │ false      │
│ 1.2.0-rc1             │ tool-F4na8SddhQ2DDM6V │ true    │ true     │     0 │ false      │
│ 1.2.0-beta1           │ tool-jNhAmK25aPJ4rWyH │ true    │ true     │     0 │ false      │
│ 1.1.9                 │ tool-DkwTW2W1cuKg1KE1 │ true    │ true     │     0 │ false      │
│ 1.1.8                 │ tool-nrKf9uxiKjCCFm8k │ true    │ true     │     0 │ false      │
│ 1.1.7                 │ tool-CFG8CQoX9fuRGjMf │ true    │ true     │     0 │ false      │
│ 1.1.6                 │ tool-BmASaTHEBG3u7wRB │ true    │ true     │     4 │ false      │
│ 1.1.5                 │ tool-tkpJshmH6uKTDw5Z │ true    │ true     │     0 │ false      │
│ 1.1.4                 │ tool-czNheBbcb9aTScf4 │ true    │ true     │     0 │ false      │
<redacted for brevity>
│ 0.6.8                 │ tool-Cbdb3tpQdLUQ2qBP │ true    │ true     │     0 │ false      │
│ 0.6.7                 │ tool-pf3wuggpDPy1ExdL │ true    │ true     │     0 │ false      │
│ 0.6.6                 │ tool-bwaiZeQSCq1yhpvf │ true    │ true     │     0 │ false      │
│ 0.6.5                 │ tool-F31bcX5VF6FNN1su │ true    │ true     │     0 │ false      │
╰───────────────────────┴───────────────────────┴─────────┴──────────┴───────┴────────────╯
```

**Search Example**

```sh
$ tfx admin terraform-version list --search 1.2.
Using config file: /Users/tstraub/.tfx.hcl
List Terraform Versions for TFE 
╭─────────────┬───────────────────────┬─────────┬──────────┬───────┬────────────╮
│ VERSION     │ ID                    │ ENABLED │ OFFICIAL │ USAGE │ DEPRECATED │
├─────────────┼───────────────────────┼─────────┼──────────┼───────┼────────────┤
│ 1.2.1       │ tool-bcVckVA7aH7b98WR │ true    │ true     │     0 │ false      │
│ 1.2.0       │ tool-7otJm3yaad55iBQY │ true    │ true     │     0 │ false      │
│ 1.2.0-rc2   │ tool-1TgFEgk2nSFfC8e4 │ true    │ true     │     0 │ false      │
│ 1.2.0-rc1   │ tool-F4na8SddhQ2DDM6V │ true    │ true     │     0 │ false      │
│ 1.2.0-beta1 │ tool-jNhAmK25aPJ4rWyH │ true    │ true     │     0 │ false      │
╰─────────────┴───────────────────────┴─────────┴──────────┴───────┴────────────╯
```

## `tfx admin terraform-version show`

Show details for a supplied Terraform Version.

**Example**

```sh
$ tfx admin terraform-version show --version 1.2.1
Using config file: /Users/tstraub/.tfx.hcl
Show Terraform Version: 1.2.1
Version: 1.2.1
ID:      tool-bcVckVA7aH7b98WR
URL:     https://releases.hashicorp.com/terraform/1.2.1/terraform_1.2.1_linux_amd64.zip
Sha:     8cf8eb7ed2d95a4213fbfd0459ab303f890e79220196d1c4aae9ecf22547302e
Enabled: true
Beta:    false
```

## `tfx admin terraform-version delete`

Delete a version of a supplied Terraform Version.

**Official Example**

:::note
When attempting to delete an "Official" Terraform Version this command will first set the the "Official" to be `false`.
:::
:::caution
Currently deleting a Terraform Version that was shipped with the installation and marked as "Official" will return the next time Terraform Enterprise is upgraded or reinstalled.
:::

```sh
$ tfx admin terraform-version delete --version 0.6.5
Using config file: /Users/tstraub/.tfx.hcl
Delete Terraform Version: 0.6.5
Forcing Terraform Version to be unofficial 
Variable Deleted: 0.6.5
Status: Success
```

## `tfx admin terraform-version create`

Create a Terraform Version.

Setting the `--disabled` flag will create the Terraform Version to be disabled (defaults to enabled).

Setting the `--official` flag will create the Terraform Version to have "Official" to be `true`.

Setting the `--beta` flag will create the Terraform Version to have "Beta" to be `true`.

Setting the `--deprecated` flag will create the Terraform Version to have "Deprecated" to be `true`.

**Example**

```sh
$ tfx admin terraform-version create --version 1.2.1 --official --url https://releases.hashicorp.com/terraform/1.2.1/terraform_1.2.1_linux_amd64.zip --sha 8cf8eb7ed2d95a4213fbfd0459ab303f890e79220196d1c4aae9ecf22547302e 
Using config file: /Users/tstraub/.tfx.hcl
Create Terraform Version: 1.2.1
Version: 1.2.1
ID:      tool-71Zae78LyPidSL84
URL:     https://releases.hashicorp.com/terraform/1.2.1/terraform_1.2.1_linux_amd64.zip
Sha:     8cf8eb7ed2d95a4213fbfd0459ab303f890e79220196d1c4aae9ecf22547302e
Enabled: true
```

## `tfx admin terraform-version create official`

Create an official Terraform Version from releases.hashicorp.com and lookup the appropriate values for `--url` and `--sha` to make adding new versions easier.

Setting the `--disabled` flag will create the Terraform Version to be disabled (defaults to enabled).

Setting the `--official` flag will create the Terraform Version to have "Official" to be `true`.

Setting the `--beta` flag will create the Terraform Version to have "Beta" to be `true`.

Setting the `--deprecated` flag will create the Terraform Version to have "Deprecated" to be `true`.

**Example**

```sh
$ tfx admin terraform-version create official --version 1.2.1
Using config file: /Users/tstraub/.tfx.hcl
Searching for official Terraform Version: 1.2.1
Terraform Version SHASUM: 8cf8eb7ed2d95a4213fbfd0459ab303f890e79220196d1c4aae9ecf22547302e
Create Terraform Version: 1.2.1
Version: 1.2.1
ID:      tool-3rPfWYPuwkgxGz4X
URL:     https://releases.hashicorp.com/terraform/1.2.1/terraform_1.2.1_linux_amd64.zip
Sha:     8cf8eb7ed2d95a4213fbfd0459ab303f890e79220196d1c4aae9ecf22547302e
Enabled: true
Beta:    false
```

## `tfx admin terraform-version disable`

Disable a Terraform Version(s), accepts comma separated list.
This command will attempt to disable all given versions even if there are failures.

Successful disables report `Disabled` in the result output.

**Basic Example**

```sh
$ tfx admin terraform-version disable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-rc2
Using config file: /Users/tstraub/.tfx.hcl
Disable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-rc2]
1.2.0-beta1: Disabled
1.2.0-rc1:   Disabled
1.2.0-rc2:   Disabled
```

**Error Example**

```sh
$ tfx admin terraform-version disable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-nope
Using config file: /Users/tstraub/.tfx.hcl
Disable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-nope]
1.2.0-beta1: Disabled
1.2.0-rc1:   Disabled
1.2.0-nope:  failed to find terraform version
```

## `tfx admin terraform-version disable all`

Disable Terraform Versions within the Terraform Enterprise install. With no filter flags, every version is targeted. Optional filter flags select a subset; only one filter flag may be used per invocation.

This command will attempt to disable all matching versions even if there are failures. Versions currently in use by workspaces cannot be disabled and are reported with `unable to disable a terraform version in use`.

| Flag | Description |
|---|---|
| `--except` | Comma-separated keep-list; disable all versions **not** in the list |
| `--before` | Disable all versions strictly before the given semver (e.g., `1.12.0` disables `1.11.x` and older) |
| `--not-in-use` | Disable only versions with `Usage == 0` |
| `--beta` | Disable only versions marked as beta |
| `--deprecated` | Disable only deprecated versions |
| `--unofficial` | Disable only unofficial versions |
| `--official` | Disable only official versions |

:::note
This command can take up to a minute to run (or longer depending on network latency).
:::

**Disable all**

```sh
$ tfx admin terraform-version disable all
Using config file: /Users/tstraub/.tfx.hcl
Disabling all Terraform versions
2.0.0:                 Disabled
1.2.1:                 Disabled
1.2.0:                 Disabled
1.1.6:                 unable to disable a terraform version in use
<redacted for brevity>
```

**Disable all except a keep-list**

```sh
$ tfx admin terraform-version disable all --except 1.12.0,1.13.0
Using config file: /Users/tstraub/.tfx.hcl
Disabling all Terraform versions except: [1.12.0 1.13.0]
2.0.0: Disabled
1.2.1: Disabled
<redacted for brevity>
1.12.0: (not targeted — kept enabled)
1.13.0: (not targeted — kept enabled)
```

**Disable all versions before a semver**

```sh
$ tfx admin terraform-version disable all --before 1.12.0
Using config file: /Users/tstraub/.tfx.hcl
Disabling all Terraform versions before 1.12.0
1.11.4: Disabled
1.11.3: Disabled
1.10.5: Disabled
<redacted for brevity>
```

Versions that cannot be parsed as semver (or are equal to or after the cutoff) are skipped.

**Disable unused versions only**

```sh
$ tfx admin terraform-version disable all --not-in-use
Using config file: /Users/tstraub/.tfx.hcl
Disabling unused Terraform versions
2.0.0: Disabled
1.2.1: Disabled
<redacted for brevity>
```

**Disable beta versions only**

```sh
$ tfx admin terraform-version disable all --beta
Using config file: /Users/tstraub/.tfx.hcl
Disabling beta Terraform versions
1.6.0-beta1: Disabled
```

**Disable deprecated versions only**

```sh
$ tfx admin terraform-version disable all --deprecated
Using config file: /Users/tstraub/.tfx.hcl
Disabling deprecated Terraform versions
```

**Disable unofficial versions only**

```sh
$ tfx admin terraform-version disable all --unofficial
Using config file: /Users/tstraub/.tfx.hcl
Disabling unofficial Terraform versions
```

**Disable official versions only**

```sh
$ tfx admin terraform-version disable all --official
Using config file: /Users/tstraub/.tfx.hcl
Disabling official Terraform versions
```

## `tfx admin terraform-version enable`

Enables a Terraform Version(s), accepts comma separated list.
This command will attempt to enable all given versions even if there are failures.

Successful enables report `Enabled` in the result output.

**Basic Example**

```sh
$ tfx admin terraform-version enable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-rc2 
Using config file: /Users/tstraub/.tfx.hcl
Enable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-rc2]
1.2.0-beta1: Enabled
1.2.0-rc1:   Enabled
```

**Error Example**

```sh
$ tfx admin terraform-version enable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-nope
Using config file: /Users/tstraub/.tfx.hcl
Enable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-nope]
1.2.0-beta1: Enabled
1.2.0-rc1:   Enabled
1.2.0-nope:  failed to find terraform version
```

## `tfx admin terraform-version enable all`

Enable Terraform Versions within the Terraform Enterprise install. With no filter flags, every version is targeted. Optional filter flags select a subset; only one filter flag may be used per invocation.

This command will attempt to enable all matching versions even if there are failures.

| Flag | Description |
|---|---|
| `--include` | Comma-separated allow-list; enable only these versions |
| `--except` | Comma-separated skip-list; enable all versions **not** in the list |
| `--beta` | Enable only versions marked as beta |
| `--unofficial` | Enable only unofficial versions |
| `--official` | Enable only official versions |

:::note
This command can take up to a minute to run (or longer depending on network latency).
:::

**Enable all**

```sh
$ tfx admin terraform-version enable all
Using config file: /Users/tstraub/.tfx.hcl
Enabling all Terraform versions
2.0.0:                 Enabled
1.2.1:                 Enabled
1.2.0:                 Enabled
<redacted for brevity>
```

**Enable only specific versions**

```sh
$ tfx admin terraform-version enable all --include 1.12.0,1.13.0
Using config file: /Users/tstraub/.tfx.hcl
Enabling Terraform versions: [1.12.0 1.13.0]
1.12.0: Enabled
1.13.0: Enabled
```

**Enable all except a skip-list**

```sh
$ tfx admin terraform-version enable all --except 1.12.0,1.13.0
Using config file: /Users/tstraub/.tfx.hcl
Enabling all Terraform versions except: [1.12.0 1.13.0]
2.0.0: Enabled
1.2.1: Enabled
<redacted for brevity>
1.12.0: (not targeted — left disabled)
1.13.0: (not targeted — left disabled)
```

**Enable beta versions only**

```sh
$ tfx admin terraform-version enable all --beta
Using config file: /Users/tstraub/.tfx.hcl
Enabling beta Terraform versions
1.6.0-beta1: Enabled
```

**Enable unofficial versions only**

```sh
$ tfx admin terraform-version enable all --unofficial
Using config file: /Users/tstraub/.tfx.hcl
Enabling unofficial Terraform versions
```

**Enable official versions only**

```sh
$ tfx admin terraform-version enable all --official
Using config file: /Users/tstraub/.tfx.hcl
Enabling official Terraform versions
```
