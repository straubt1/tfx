---
title: Private Registry Provider Commands
---

The ability to manage Providers within an Organization was added to Terraform Enterprise in release v202206-1.
These commands make the management of these providers via the API (the only way to manage said providers) easier.

To copy a provider from the public Terraform Registry into a private registry, see [Syncing Providers to a Private Registry](/guides/sync-providers/).

> Note: Create, show, list, and delete operate on **private** registry providers. `tfx registry provider download` fetches artifacts from the **public** Terraform Registry so they can be staged and then uploaded with the existing create commands.

There are several "resources" needed to create a Provider in the Registry that have a dependency hierarchy.

``` mermaid
classDiagram
  Provider --|> ProviderVersion
  ProviderVersion --|> ProviderVersionPlatform
  class Provider{
    +String Name, Example: "aws"
  }
  class ProviderVersion{
    +String Version, Example: "1.0.0"
    +String SHASUMS
    +String SHASUMSSIG
    +String GPG Key Id
  }
  class ProviderVersionPlatform{
    +String OS, Example: "linux"
    +String ARCH, Example: "amd64"
    +String Binary
  }
```

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

## `tfx registry provider download`

Download a provider from the public Terraform Registry and stage the files locally for a later private-registry upload.

This command does **not** require a TFE/HCP Terraform token. It uses the public registry protocol (`registry.terraform.io`) and follows the same URLs Terraform uses during `terraform init`.

`--directory` is the base path (default `./providers`). Files are always written to `<directory>/<namespace>/<name>/<version>/`. The public-registry namespace (who published the provider on `registry.terraform.io`) is a folder in that path. Official providers use `--namespace hashicorp` (the default). Third-party providers use their publisher, e.g. `--namespace chainguard-dev`.

Default platforms: `linux_amd64`, `darwin_arm64`, `darwin_amd64`, `windows_amd64`. Use `--platforms` to choose a subset, or `--all-platforms` to fetch every published zip.

The GPG public key is always saved as `<KEYID>.asc`, including for the `hashicorp` namespace. Official HashiCorp providers are signed with HashiCorp's public GPG keys, which are publicly available and pre-installed on Terraform Enterprise v202309-1 and newer. Third-party providers use their own keys, which are **not** pre-installed; `tfx registry provider version create --directory` uploads them to the private registry (the TFE/HCP organization) when missing.

After upload, Terraform configs must source the provider as `<hostname>/<organization>/azurerm` — not `hashicorp/azurerm`.

Unlike `tfx registry module version download` (from the private registry), this command downloads **from the public registry**.

**Example:**

```sh
$ tfx registry provider download --name azurerm --version 5.0.0
Download Provider from Public Registry
Provider: hashicorp/azurerm 5.0.0
Platforms: linux_amd64, darwin_arm64, darwin_amd64, windows_amd64
Downloading provider artifacts...
Namespace:     hashicorp
Name:          azurerm
Version:       5.0.0
Directory:     ~/example/providers/hashicorp/azurerm/5.0.0
GPG Key ID:    <hashicorp signing key>
GPG Public Key: ~/example/providers/hashicorp/azurerm/5.0.0/<KEYID>.asc
SHA256SUMS:    ~/example/providers/hashicorp/azurerm/5.0.0/terraform-provider-azurerm_5.0.0_SHA256SUMS
SHA256SUMS.sig: ~/example/providers/hashicorp/azurerm/5.0.0/terraform-provider-azurerm_5.0.0_SHA256SUMS.72D7468F.sig
```

Staged layout:

```
./providers/hashicorp/azurerm/5.0.0/
  <KEYID>.asc
  terraform-provider-azurerm_5.0.0_SHA256SUMS
  terraform-provider-azurerm_5.0.0_SHA256SUMS.72D7468F.sig
  terraform-provider-azurerm_5.0.0_linux_amd64.zip
  terraform-provider-azurerm_5.0.0_darwin_arm64.zip
  terraform-provider-azurerm_5.0.0_darwin_amd64.zip
  terraform-provider-azurerm_5.0.0_windows_amd64.zip
```

Then upload with:

```sh
tfx registry provider version create --directory ./providers/hashicorp/azurerm/5.0.0
```

The provider is created automatically if it does not already exist. `--directory` uploads the version checksums and every zip in that folder.

```sh
$ tfx registry provider download --name azurerm --version 5.0.0 --platforms linux_amd64 --directory ./providers
```

### Third-party providers

Third-party publishers use their own GPG keys. Download with `--namespace` set to the public-registry publisher (not `hashicorp`):

```sh
tfx registry provider download --namespace chainguard-dev --name cosign --version 0.4.16
tfx registry provider version create --directory ./providers/chainguard-dev/cosign/0.4.16
```

`version create --directory` reads the public-registry namespace from the staged folder path (`<namespace>/<name>/<version>`). When it is not `hashicorp`, the command checks whether that provider's GPG key already exists in the private registry (the TFE/HCP organization) and creates it from the staged `.asc` file if missing.

A provider is in the HashiCorp namespace only when that path namespace equals `hashicorp` (case-insensitive). The provider name and the GPG key ID are not used for this decision.

## `tfx registry provider version list`

List Versions for a Provider in the Registry.

## `tfx registry provider version platform list`

List Platforms for a Provider Version in the Registry.

## `tfx registry provider version create`

Create a Version for a Provider in the Registry.

### From a staged download directory

After `tfx registry provider download`, pass the version folder. Namespace, name, version, GPG key id, SHA256SUMS, signature, and platform zips are inferred. The provider is created if it does not exist. For third-party providers (public-registry namespace is not `hashicorp`), the GPG public key is uploaded to the organization if it is not already present. `--key-id` overrides the inferred GPG key. Platform zips upload in parallel; `--concurrency` defaults to 4.

```sh
tfx registry provider version create --directory ./providers/hashicorp/azurerm/5.0.0
tfx registry provider version create --directory ./providers/hashicorp/aws/6.60.0 --concurrency 4
```

### From explicit files

`--shasums` is the path to the SHA256SUMS file. `--shasums-sig` is the path to the signature file. This path creates the version only; upload platforms separately with `tfx registry provider version platform create`.

**SHAMSUM File Content Example:**

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

```sh
$ tfx registry provider version create --name random --version 4.3.0 --key-id <gpg-key-id> --shasums ./terraform-provider-random_3.1.0_SHA256SUMS --shasums-sig=./terraform-provider-random_3.1.0_SHA256SUMS.sig
Using config file: /Users/tstraub/.tfx.hcl
Create Provider Version in Registry for Organization: firefly
Provider Name: random
Uploading shasums and sig 
~/hashicorp-services.github.com/pmr-providers-guide/providers/random/3.1.0/terraform-provider-random_3.1.0_SHA256SUMS ~/hashicorp-services.github.com/pmr-providers-guide/providers/random/3.1.0/terraform-provider-random_3.1.0_SHA256SUMS.sig 2022-08-20T18:27:13.859Z
Provider Version Created 
Name:    random
ID:      provver-ujrSC8txwA62a1uz
Version: 4.3.0
Created: 2022-08-20T18:27:13.859Z
```

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

Delete a Version for a Provider in the Registry.

**Example**

```sh
$ tfx registry provider version delete --name aws --version 4.3.0
Using config file: /Users/tstraub/.tfx.hcl
Delete Provider Version in Registry for Organization: firefly
Provider Version Deleted: aws
Status: Success
```

## `tfx registry provider version platform create`

Create a Platform Version for a Provider in the Registry.

**Example**

```sh
$ tfx registry provider version platform create --name aws --version 4.3.0 --os darwin --arch amd64 -f ./terraform-provider-aws_4.3.0_darwin_amd64.zip
Using config file: /Users/tstraub/.tfx.hcl
Create Provider Platform in Registry for Organization: firefly
Hashing Provider File 
Building Provider Filename terraform-provider-random_4.3.0_darwin_amd64.zip
Uploading Provider Version Platform... 
Provider Platform Created: 
ID:      provpltfrm-mCPKwkHwLyvhPchS
OS:      darwin
Arch:    amd64
```

## `tfx registry provider version platform show`

Show details of a Platform Version for a Provider in the Registry.

**Example**

```sh
$ tfx registry provider version platform show --name aws --version 4.3.0 --os darwin --arch amd64
Using config file: /Users/tstraub/.tfx.hcl
Show Provider Platform in Registry for Organization: firefly
Name:     aws
ID:       provpltfrm-mCPKwkHwLyvhPchS
Version:  4.3.0
OS:       darwin
ARCH:     amd64
Filename: terraform-provider-aws_4.3.0_darwin_amd64.zip
Shasum:   e31c31d00f42ea2dbaab1ad4c245da5cfff63e28399b5a5795b5e6a826c6c8af
```

## `tfx registry provider version platform delete`

Delete a Platform Version for a Provider in the Registry.

**Example**

```sh
$ tfx registry provider version platform delete --name aws --version 4.3.0 --os darwin --arch amd64
Using config file: /Users/tstraub/.tfx.hcl
Delete Provider Platform in Registry for Organization: firefly
Provider Version Deleted: aws
Status: Success
```
