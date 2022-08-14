# Terraform Version Commands

Managing Terraform Versions in a Terraform Enterprise install.

!!! note ""
    These commands will only work with Terraform Enterprise

!!! note ""
    All commands below can be used with a `tfv` alias.

## `tfx admin terraform-version list`

List all Terraform Versions for a Terraform Enterprise install

## `tfx admin terraform-version show`

Show version details for a supplied Terraform Version or Version Id

## `tfx admin terraform-version delete`

Delete a version of a supplied Terraform Version or Version Id

## `tfx admin terraform-version create`

Create a Terraform Version

## `tfx admin terraform-version create official`

Create an official Terraform Version from releases.hashicorp.com

## `tfx admin terraform-version disable`

Disable a Terraform Version(s), accepts comma separated list

## `tfx admin terraform-version disable all`

Disable all Terraform Versions

## `tfx admin terraform-version enable`

Enables a Terraform Version(s), accepts comma separated list

## `tfx admin terraform-version enable all`

Enables all Terraform Versions
