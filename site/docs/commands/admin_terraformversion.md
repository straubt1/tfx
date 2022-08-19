# Terraform Version Commands

Managing Terraform Versions in a Terraform Enterprise install.

!!! note ""
    These commands will only work with Terraform Enterprise

!!! note ""
    All commands below can be used with a `tfv` alias.

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

!!! note ""
    When attempting to delete an "Official" Terraform Version this command will first set the the "Official" to be `false`.
!!! warning
    Currently deleting a Terraform Version that was shipped with the installation and marked as "Official" will return the next time Terraform Enterprise is upgraded or reinstalled.

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
This command will attempt to enable all given versions even if there are failures.

**Basic Example**

```sh
$ tfx admin terraform-version disable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-rc2
Using config file: /Users/tstraub/.tfx.hcl
Disable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-rc2]
1.2.0-beta1: false
1.2.0-rc1:   false
1.2.0-rc2:   false
```

**Error Example**

```sh
$ tfx admin terraform-version disable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-nope
Using config file: /Users/tstraub/.tfx.hcl
Disable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-nope]
1.2.0-beta1: false
1.2.0-rc1:   false
1.2.0-nope:  failed to find terraform version
```

## `tfx admin terraform-version disable all`

Disable all Terraform Versions within the Terraform Enterprise install.
This command will attempt to enable all given versions even if there are failures.

!!! note ""
    This command can take up to a minute to run (or longer depending on network latency).

**Example**

```sh
$ tfx admin terraform-version disable all
Using config file: /Users/tstraub/.tfx.hcl
Disable All Terraform Versions 
2.0.0:                 false
1.2.1:                 false
1.2.0:                 false
1.2.0-rc2:             false
1.2.0-rc1:             false
1.2.0-beta1:           false
1.1.9:                 false
<redacted for brevity>
0.6.12:                false
0.6.11:                false
0.6.10:                false
0.6.9:                 false
0.6.8:                 false
0.6.7:                 false
0.6.6:                 false
```

## `tfx admin terraform-version enable`

Enables a Terraform Version(s), accepts comma separated list.
This command will attempt to enable all given versions even if there are failures.

**Basic Example**

```sh
$ tfx admin terraform-version enable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-rc2 
Using config file: /Users/tstraub/.tfx.hcl
Enable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-rc2]
1.2.0-beta1: true
1.2.0-rc1:   true
```

**Error Example**

```sh
$ tfx admin terraform-version enable --versions 1.2.0-beta1,1.2.0-rc1,1.2.0-nope
Using config file: /Users/tstraub/.tfx.hcl
Enable Terraform Versions: [1.2.0-beta1 1.2.0-rc1 1.2.0-nope]
1.2.0-beta1: true
1.2.0-rc1:   true
1.2.0-nope:  failed to find terraform version
```

## `tfx admin terraform-version enable all`

Disable all Terraform Versions within the Terraform Enterprise install.
This command will attempt to enable all given versions even if there are failures.

**Example**

```sh
$ tfx admin terraform-version enable all
Using config file: /Users/tstraub/.tfx.hcl
Enable All Terraform Versions 
2.0.0:                 true
1.2.1:                 true
1.2.0:                 true
1.2.0-rc2:             true
1.2.0-rc1:             true
1.2.0-beta1:           true
1.1.9:                 true
<redacted for brevity>
0.6.12:                true
0.6.11:                true
0.6.10:                true
0.6.9:                 true
0.6.8:                 true
0.6.7:                 true
0.6.6:                 true
```
