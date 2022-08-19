# Private Registry Provider Commands

The ability to manage Providers within an Organization was added to Terraform Enterprise in release v202206-1.
These commands make the management of these providers via the API (the only way to manage said providers) easier.

> Note: These commands only work with the "private" Providers.

There are several "resources" needed to create a Provider in the Registry that have a dependency hierarchy.

- "Provider" 
  - (Name, Example: "aws")
  - "Provider Version" 
    - (Version, Example: "4.0.1") 
    - (SHASUMS & SHASUMSSIG files)
    - (GPG Key ID)
    - "Provider Version Platform"
      - (OS, Example: Linux)
      - (Arch, Example: amd64) 

## `tfx registry provider list`

List Providers in the Registry.

**Example:**

```sh
$ tfx registry provider list
Using config file: /Users/tstraub/.tfx.hcl
List Providers in Registry for Organization: firefly
╭─────────┬──────────┬───────────────────────┬──────────────────────────╮
│ NAME    │ REGISTRY │ ID                    │ PUBLISHED                │
├─────────┼──────────┼───────────────────────┼──────────────────────────┤
│ aws     │ private  │ prov-xT3mHHwFEe9BA62L │ 2022-07-15T18:30:33.218Z │
│ azurerm │ private  │ prov-zhufbe1gdvHxyrzV │ 2022-07-15T18:41:04.675Z │
│ random  │ private  │ prov-A8tqWgWykT3ecb1h │ 2022-07-15T18:43:37.771Z │
╰─────────┴──────────┴───────────────────────┴──────────────────────────╯
```

## `tfx registry provider create`

Create a Provider in the Registry.

**Example:**

```sh
$ tfx registry provider create --name google
Using config file: /Users/tstraub/.tfx.hcl
Create Provider in Registry for Organization: firefly
Provider Created: google
ID:        prov-tAaS9tEKZFMTr53f
Namespace: firefly
Created:   2022-08-13T17:34:38.067Z
```

## `tfx registry provider show`

Show details of a Provider in the Registry.

**Example:**

```sh
$ tfx registry provider show --name google
Using config file: /Users/tstraub/.tfx.hcl
Show Provider in Registry for Organization: firefly
Name:      google
ID:        prov-tAaS9tEKZFMTr53f
Namespace: firefly
Created:   2022-08-13T17:34:38.067Z
```

## `tfx registry provider delete`

Delete a Provider in the Registry.

**Example:**

```sh
$ tfx registry provider delete --name google
Using config file: /Users/tstraub/.tfx.hcl
Delete Provider in Registry for Organization: firefly
Provider Deleted: google
Status: Success
```

## `tfx registry provider version list`

List Versions for a Provider in the Registry.

## `tfx registry provider version platform list`

List Platforms for a Provider Version in the Registry.

## `tfx registry provider version create`

Create a Version for a Provider in the Registry.

`--shasums` Is required to be set to the path to shasums file. This file contains all the SHASUMS for each provider version platform you wish to upload.

**Example:**

```
e31c31d00f42ea2dbaab1ad4c245da5cfff63e28399b5a5795b5e6a826c6c8af  terraform-provider-aws_4.3.0_darwin_amd64.zip
de166ecfeed70f570cea72ec094f00c2f997496b3226fa08518e7cd4a73884e1  terraform-provider-aws_4.3.0_darwin_arm64.zip
f93725afd8410194ede51d83505327aa1ae6a9b4280cf31db649c62c7dc203ae  terraform-provider-aws_4.3.0_freebsd_386.zip
087c67e5429f343a164221c05a83f152322f411e7394f8a39ed81a75982af1f2  terraform-provider-aws_4.3.0_freebsd_amd64.zip
2e852a1b107e5324524874e1cd98bcf3a69284b4fe04750aa373054177c54214  terraform-provider-aws_4.3.0_freebsd_arm.zip
4b9a54b5895f945827832e6ddd16ff107301fedf47acbd83d17d4e18bbf10bb1  terraform-provider-aws_4.3.0_linux_386.zip
64dfc02bc85f5df2f51ff942fc78d72fcd0db17b0f53e1fae380e58adbd239b3  terraform-provider-aws_4.3.0_linux_amd64.zip
c51f5b238af37c63e9033a12fd7fedc87c03eb966f5f5c7786eb6246e8bf3071  terraform-provider-aws_4.3.0_linux_arm.zip
d0df94d3112a25de609dfb55c5e3b0d119dea519a2bdd8099e64a8d63f22b683  terraform-provider-aws_4.3.0_linux_arm64.zip
90048d87ff3071a4356cf91916b46a7ec69ba55bcba5765b598d3fe545d4c6ca  terraform-provider-aws_4.3.0_windows_386.zip
766f9aef619cfd23e924aee523791acccd30b6d8f1cc0ed1a7b5c953bf8c5392  terraform-provider-aws_4.3.0_windows_amd64.zip
```

`--shasumssig` Is required to be set to the path to shasums signature file. This file is a binary

## `tfx registry provider version show`

Show details a Version for a Provider in the Registry.

**Example:**

```sh
$ tfx registry provider version show --name aws --version 4.3.0
Using config file: /Users/tstraub/.tfx.hcl
Show Provider Version in Registry for Organization: firefly
Name:                 aws
Version:              4.3.0
ID:                   provver-pmXF2YXLARN8ZpYp
Shasums Uploaded:     true
Shasums Sig Uploaded: true
Shasums:              
e31c31d00f42ea2dbaab1ad4c245da5cfff63e28399b5a5795b5e6a826c6c8af  terraform-provider-aws_4.3.0_darwin_amd64.zip
de166ecfeed70f570cea72ec094f00c2f997496b3226fa08518e7cd4a73884e1  terraform-provider-aws_4.3.0_darwin_arm64.zip
f93725afd8410194ede51d83505327aa1ae6a9b4280cf31db649c62c7dc203ae  terraform-provider-aws_4.3.0_freebsd_386.zip
087c67e5429f343a164221c05a83f152322f411e7394f8a39ed81a75982af1f2  terraform-provider-aws_4.3.0_freebsd_amd64.zip
2e852a1b107e5324524874e1cd98bcf3a69284b4fe04750aa373054177c54214  terraform-provider-aws_4.3.0_freebsd_arm.zip
4b9a54b5895f945827832e6ddd16ff107301fedf47acbd83d17d4e18bbf10bb1  terraform-provider-aws_4.3.0_linux_386.zip
64dfc02bc85f5df2f51ff942fc78d72fcd0db17b0f53e1fae380e58adbd239b3  terraform-provider-aws_4.3.0_linux_amd64.zip
c51f5b238af37c63e9033a12fd7fedc87c03eb966f5f5c7786eb6246e8bf3071  terraform-provider-aws_4.3.0_linux_arm.zip
d0df94d3112a25de609dfb55c5e3b0d119dea519a2bdd8099e64a8d63f22b683  terraform-provider-aws_4.3.0_linux_arm64.zip
90048d87ff3071a4356cf91916b46a7ec69ba55bcba5765b598d3fe545d4c6ca  terraform-provider-aws_4.3.0_windows_386.zip
766f9aef619cfd23e924aee523791acccd30b6d8f1cc0ed1a7b5c953bf8c5392  terraform-provider-aws_4.3.0_windows_amd64.zip
```

## `tfx registry provider version delete`

Delete a Version for a Provider in the Registry

## `tfx registry provider version platform create`

Create a Platform Version for a Provider in the Registry

## `tfx registry provider version platform show`

Show details of a Platform Version for a Provider in the Registry

## `tfx registry provider version platform delete`

Delete a Platform Version for a Provider in the Registry
