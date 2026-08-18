---
title: GPG Commands
---

Managing GPG Keys in an Organization.

Currently these are only used for creating Providers in the Private Registry.

Official HashiCorp providers are signed with HashiCorp's public GPG keys, which are publicly available and pre-installed on Terraform Enterprise v202309-1 and newer. Third-party providers use their own keys; upload those with `tfx admin gpg create` (or let `tfx registry provider version create --directory` do it from a staged `.asc`).

## `tfx admin gpg list`

List GPG Keys of an Organization.

**Example**

```sh
$ tfx admin gpg list   
Using config file: /Users/tstraub/.tfx.hcl
List GPG Keys for Organization: firefly
╭──────────────────┬───────────┬───────────────────────┬───────────────────────╮
│ KEY ID           │ NAMESPACE │ UPDATED AT            │ CREATED AT            │
├──────────────────┼───────────┼───────────────────────┼───────────────────────┤
│ 37AD5AEF6A5D6D5C │ firefly   │ Sat Jul  9 16:26 2022 │ Sat Jul  9 16:26 2022 │
╰──────────────────┴───────────┴───────────────────────┴───────────────────────╯
```

## `tfx admin gpg create`

Create a GPG Key for an Organization using the public key contents.

**Example**

```sh
$ tfx admin gpg create --namespace firefly --public-key ./provider.gpg.public 
Using config file: /Users/tstraub/.tfx.hcl
Create GPG Key for Organization: firefly
GPG Key Created 
KeyID:      37AD5AEF6A5D6D5C
Created:    Fri Aug 19 21:48 2022
Updated:    Fri Aug 19 21:48 2022
AsciiArmor: 
-----BEGIN PGP PUBLIC KEY BLOCK-----
...
-----END PGP PUBLIC KEY BLOCK-----
```

## `tfx admin gpg show`

Show details of a GPG Key for an Organization.

**Example**

```sh
$ tfx admin gpg show --namespace firefly --id 37AD5AEF6A5D6D5C
Using config file: /Users/tstraub/.tfx.hcl
Show a GPG Key for Organization: firefly
GPG Key Found 
KeyID:      37AD5AEF6A5D6D5C
Created:    Sat Jul  9 16:26 2022
Updated:    Sat Jul  9 16:26 2022
AsciiArmor: 
-----BEGIN PGP PUBLIC KEY BLOCK-----
...
-----END PGP PUBLIC KEY BLOCK-----
```

## `tfx admin gpg delete`

Delete GPG Key for an Organization.

:::caution
The API will allow you to delete a GPG that is in use, caution advised.
:::

**Example**

```sh
$ tfx admin gpg delete --namespace firefly --id 37AD5AEF6A5D6D5C
Using config file: /Users/tstraub/.tfx.hcl
Delete GPG Key for Organization: firefly
GPG Key Deleted 
Status: Success
```
