---
title: Syncing Providers to a Private Registry
description: Download providers from the public Terraform Registry and publish them to a private registry on HCP Terraform or Terraform Enterprise.
---

Copy a provider from the public Terraform Registry into the private registry of an HCP Terraform or Terraform Enterprise organization. TFx does this in two commands: download artifacts locally, then publishes them.

This works the same way on **HCP Terraform** and **Terraform Enterprise**. Only the hostname in your profile and in `required_providers` source strings changes.

For flag-level detail, see [Private Registry Provider Commands](/commands/registry_provider/) and [GPG Commands](/commands/admin_gpg/).

## Why sync providers?

- **Allowlist what teams can use.** Publish only the providers and versions you want in the organization. Workspaces then consume your private registry instead of the entire public catalog.
- **Survive public-registry outages or limited egress.** Once the binaries live in your private registry, `terraform init` does not depend on `registry.terraform.io` being reachable.
- **Treat providers as executables.** Provider zips are binaries you run. Copying them through your registry is a natural place to scan, review, and pin versions before they reach workspaces.

:::note
After publish, Terraform configs must source the provider as `<hostname>/<organization>/<name>` — not `hashicorp/<name>`. See [Use the provider](#use-the-provider).
:::

## Step by step

Download does **not** require a TFE or HCP Terraform token. It communicates with `registry.terraform.io`. Upload uses your TFx profile.

### 1. Authenticate

Target the organization that owns the private registry (`tfx login`, or `--profile`).

### 2. Download from the public registry

```sh
tfx registry provider download --name azurerm --version 5.0.0
```

Files are written to `<directory>/<namespace>/<name>/<version>/`. `--directory` defaults to `./providers`. Official providers use `--namespace hashicorp` (the default).

Default platforms: `linux_amd64`, `darwin_arm64`, `darwin_amd64`, `windows_amd64`. Use `--platforms` for a subset, or `--all-platforms` to fetch every published zip.

```sh
tfx registry provider download --name azurerm --version 5.0.0 --platforms linux_amd64
tfx registry provider download --name azurerm --version 5.0.0 --all-platforms
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

### 3. Upload the staged folder

`tfx registry provider version create --directory` infers namespace, name, version, GPG key, checksums, and platform zips. It creates the provider if it does not already exist — no separate `tfx registry provider create` is required.

```sh
tfx registry provider version create --directory ./providers/hashicorp/azurerm/5.0.0
```

Platform archive files are uploaded in parallel. `--concurrency` defaults to 4.

### 4. Use the provider

Point Terraform at the private registry. The source is `<hostname>/<organization>/<name>`:

```hcl
terraform {
  required_providers {
    azurerm = {
      source  = "app.terraform.io/my-org/azurerm" # or tfe.example.com/my-org/azurerm
      version = "5.0.0"
    }
  }
}
```

## Third-party providers

Third-party publishers use their own GPG keys. Set `--namespace` to the public-registry publisher (not `hashicorp`):

```sh
tfx registry provider download --namespace chainguard-dev --name cosign --version 0.4.16
tfx registry provider version create --directory ./providers/chainguard-dev/cosign/0.4.16
```

The public-registry namespace is a folder in the staged path. It is not the private-registry organization. A provider is treated as HashiCorp only when that path namespace equals `hashicorp` (case-insensitive).

## GPG keys

Every download writes a `<KEYID>.asc` next to the SHA256SUMS files, including HashiCorp providers.

- **Official HashiCorp providers** (`--namespace hashicorp`) are signed with HashiCorp's public GPG keys. Those keys are pre-installed on Terraform Enterprise v202309-1 and newer, and on HCP Terraform. Directory-mode create does **not** upload a HashiCorp key.
- **Third-party providers** use their own keys, which are **not** pre-installed. `version create --directory` checks the organization and creates the key from the staged `.asc` when it is missing.

The GPG namespace is the **TFE or HCP Terraform organization**, not the public-registry publisher.

You can still manage keys yourself with [`tfx admin gpg`](/commands/admin_gpg/) (`list`, `create`, `show`, `delete`). Pass `--key-id` on `version create` to override the key inferred from the staged folder.

:::caution
The API will allow you to delete a GPG key that is still in use.
:::
