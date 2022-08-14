# GPG Commands

Managing GPG Keys in an Organization.

Currently these are only used for creating Providers in the Private Registry.

## `tfx admin gpg list`

List GPG Keys of an Organization.

## `tfx admin gpg create`

Create a GPG Key for an Organization using the public key contents.

## `tfx admin gpg show`

Show details of a GPG Key for an Organization.

## `tfx admin gpg delete`

Delete GPG Key for an Organization.

!!! warning ""
    The API will allow you to delete a GPG that is in use, caution advised.
