# Command Style Guide

A good CLI has a consistent and concise feel.

The goal of this document is to outline the general user experience when dealing with the default output. 

!!! Note ""
  When using the `--json` flag on a command, the output will simply be JSON and this guide does not apply.

## Colors

- User Provided: Green
  - These values are provided by the user, either through a configuration file or a command flag.
- Calculated: Yellow
  - These values are calculated, but could be derived from user provided values.
- Remote: Blue
  - These values are read from the external system, such as TFE or TFC.
- Error: Red
  - It is the way Rick Sanchez would want it...

## Command and Flag Naming

Command and Flag naming will adhere to the following standards:

- When a command or flag spans multiple words, the words will be separated with a dash "-"
  - Examples:
    - "Key ID" -> "key-id"
    - "Configuration Version" -> "configuration-version"
- When a name is required as an input for a resource that is directly associated with the command, "name" will be used.
  - Examples:
    - "tfx workspace show --name tt-workspace"
- When a name is required as an input for a resource that is indirectly associated with the command, "{resource}-name" will be used.
  - Examples:
    - "tfx workspace run list --workspace-name tt-workspace"

## List Commands

`list` Commands will adhered to the following standards:

- For resources that common workflows typically only require recent results and can have many results:
  - The output items will default to 10.
  - The flag `--max-items` will allow a specific maximum.
  - The flag `--all` will return all results regardless of other flag values.
