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
- When a name or id is required as an input for a resource that is directly associated with the command, "name" or "id" will be used.
    - Examples:
        - "tfx workspace show --name tt-workspace"
        - "tfx workspace configuration-version show --id cv-e83GeSpjVKXuUGmU"
- When a name or id is required as an input for a resource that is indirectly associated with the command, "{resource}-name" or "{resource}-id" will be used.
    - Examples:
        - "tfx workspace run list --workspace-name tt-workspace"
        - "tfx workspace run create --workspace-name tt-workspace --configuration-version-id cv-e83GeSpjVKXuUGmU"

## List Commands

`list` Commands will adhere to the following standards:

- For resources that common workflows typically only require recent results and can have many results:
  - The output items will default to 10.
  - The flag `--max-items` will allow a specific maximum.
  - The flag `--all` will return all results regardless of other flag values.

## Download Commands

`download` Commands will adhere to the following standards:

- Give the option to pass in a Directory.
- When no Directory is given, default to creating a temporary directory