# Workspace Configuration Version Commands

Managing Workspace Configuration Versions in a Workspace.

!!! note ""
    All commands below can be used with a `cv` alias.

!!! note ""
    Deleting a Configuration Version is not possible at this time.

## `tfx workspace configuration-version list`

List all Configuration Versions for a supplied Workspace.

`--max-items` defaults to 10, setting this to a higher number will retrieve more items.

**Example**

```sh
$ tfx workspace configuration-version list --workspace-name tfx-test
Using config file: /Users/tstraub/.tfx.hcl
List Configuration Versions for Workspace: tfx-test
╭─────────────────────┬─────────────┬──────────┬──────┬────────┬────────┬─────────╮
│ ID                  │ SPECULATIVE │ STATUS   │ REPO │ BRANCH │ COMMIT │ MESSAGE │
├─────────────────────┼─────────────┼──────────┼──────┼────────┼────────┼─────────┤
│ cv-TSYCiigC5yJ5BxXw │ true        │ uploaded │      │        │        │         │
│ cv-K4EwbnFK4MGG2Qs1 │ false       │ uploaded │      │        │        │         │
│ cv-LdcmfSz6PAswZo5L │ true        │ uploaded │      │        │        │         │
│ cv-VYikVwjgfHNnUYfr │ false       │ uploaded │      │        │        │         │
│ cv-LcVa2hMVZg1nAX6R │ false       │ uploaded │      │        │        │         │
│ cv-q9yhRwv73u6UFJdq │ false       │ uploaded │      │        │        │         │
│ cv-pn7T5L8J58FV5PSZ │ false       │ uploaded │      │        │        │         │
│ cv-p8XXa5rcph3W1MoF │ false       │ uploaded │      │        │        │         │
│ cv-Ud4YHRbViJyB9sqc │ false       │ uploaded │      │        │        │         │
│ cv-2C5nthzX7mAPMsue │ false       │ uploaded │      │        │        │         │
╰─────────────────────┴─────────────┴──────────┴──────┴────────┴────────┴─────────╯
```

## `tfx workspace configuration-version create`

Create a Configuration Version for a supplied Workspace.

**Example**

```sh
$ tfx workspace configuration-version create --workspace-name tt-workspace --directory ./tt-workspace-code/
Using config file: /Users/tstraub/.tfx.hcl
Create Configuration Version for Workspace: tt-workspace
Code Directory: ./tt-workspace-code/
Upload code to Configuration Version... 
Configuration Version Created 
ID:          cv-e83GeSpjVKXuUGmU
Speculative: false
```

## `tfx workspace configuration-version show`

Show Configuration Version details for a supplied Configuration.

**Example**

```sh
$ tfx workspace configuration-version show --id cv-K4EwbnFK4MGG2Qs1
Using config file: /Users/tstraub/.tfx.hcl
Show Configuration Version for Workspace from Id: cv-K4EwbnFK4MGG2Qs1
ID:          cv-K4EwbnFK4MGG2Qs1
Status:      uploaded
Speculative: false
```

## `tfx workspace configuration-version download`

Download the Terraform code in a Configuration Version.

**Temp Folder Example**

```sh
$ tfx workspace configuration-version download --id cv-K4EwbnFK4MGG2Qs1
Using config file: /Users/tstraub/.tfx.hcl
Downloading Configuration Version from Id: cv-K4EwbnFK4MGG2Qs1
Directory not supplied, creating a temp directory 
Configuration Version Found, download started... 
Status:    Success
Directory: /var/folders/99/srh_6psj6g5520gwyv8v3nbw0000gn/T/slug3146610843/cv-K4EwbnFK4MGG2Qs1
```

**Specific Folder Example**

```sh
TODO
```
