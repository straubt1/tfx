# Command Style Guide

A good CLI has a consistent and concise feel.

The goal of this document is to outline the general user experience when dealing with the default output. 

## Colors

!!! Note ""
    When using the `--json` flag on a command these colors will not apply.

- User Provided: Green
    - These values are provided by the user, either through a configuration file or a command flag.
- Calculated: Yellow
    - These values are calculated, but could be derived from user provided values.
- Remote: Blue
    - These values are read from the external system, such as TFE or TFC.
- Error: Red
    - It is the way Rick Sanchez would want it...

## Output Spacing

For all outputs where the Value (in the Key/Value message) is a primitive, the Key will be left justified, the Value will be spaced to be even with other messages returned by the command.

Example:

```
ID:                   ws-VxepewkunumUbR9V
Terraform Version:    1.0.0
Execution Mode:       remote
Locked:               false
Global State Sharing: false
Current Run Id:       run-tNGxao7zMos5YrY1
```

For outputs where the Value is a List, the Key will be left justified, the Values of the List will be on a new line with a two space left-padding.

Example:

```
Some List:           
  item1
  item2
  item3
```

For outputs where the Value is a Map, the Key will be left justified, the Values of the Map will be on a new line with a two space left-padding for the inner key and be spaced to be even with other inner keys in the message.

Example:

```
Some Map:            
  item3:           5
  item4Long:       5.3
  item1:           string_value 1
  item2ReallyLong: true
```

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