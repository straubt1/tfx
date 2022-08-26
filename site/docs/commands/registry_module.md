# Private Registry Module Commands

Managing modules in the Private Registry.

Namespace will be the Organization Name.

There are several "resources" needed to create a Module in the Registry that have a dependency hierarchy.

``` mermaid
classDiagram
  Module --|> ModuleVersion
  class Module{
    +String Name, Example: "my-awesome-app"
    +String Provider, Example: "aws"
  }
  class ModuleVersion{
    +String Version, Example: "1.0.0"
  }
```

## `tfx registry module list`

List all modules in the Private Registry.

**Example**

```sh
$ tfx registry module list
Using config file: /Users/tstraub/.tfx.hcl
List Modules for Organization: firefly
╭────────────┬──────────┬──────────────────────┬────────────────┬──────────────────────────┬──────────╮
│ NAME       │ PROVIDER │ ID                   │ STATUS         │ PUBLISHED                │ VERSIONS │
├────────────┼──────────┼──────────────────────┼────────────────┼──────────────────────────┼──────────┤
│ tt-module  │ random   │ mod-UsfngMNiyHpG76Yh │ setup_complete │ 2021-05-15T19:24:23.256Z │        2 │
│ ss-module  │ aws      │ mod-2qugFvqojrAQNKUs │ setup_complete │ 2021-05-12T20:57:04.494Z │        2 │
│ rr-module  │ aws      │ mod-x61ktgeM4eLA6zPX │ setup_complete │ 2021-05-16T20:52:04.974Z │        1 │
╰────────────┴──────────┴──────────────────────┴────────────────┴──────────────────────────┴──────────╯
```

## `tfx registry module create`

Create a module in the Private Registry.

**Example**

```sh
$ tfx registry module create --name tt-module --provider random                
Using config file: /Users/tstraub/.tfx.hcl
Create Module for Organization: firefly
Module Created: tt-module
ID:        mod-3Gi7SGzjESjcUezn
Namespace: firefly
Created:   2022-08-19T22:08:21.861Z
```

## `tfx registry module show`

Show module details for a module in the Private Registry.

**Example**

```sh
$ tfx registry module show --name tt-module --provider random
Using config file: /Users/tstraub/.tfx.hcl
Show Module for Organization: firefly
ID:             mod-UsfngMNiyHpG76Yh
Status:         setup_complete
Created:        2021-05-15T19:23:19.914Z
Updated:        2021-05-15T19:24:23.256Z
Versions:       2
Latest Version: 0.0.1
```

## `tfx registry module delete`

Delete a module in the Private Registry.

**Example**

```sh
$ tfx registry module delete --name tt-module --provider random         
Using config file: /Users/tstraub/.tfx.hcl
Delete Module for Organization: firefly
Module Deleted: tt-module
Status: Success
```

## `tfx registry module version list`

Show versions for a module in the Private Registry.

**Example**

```sh
$ tfx registry module version list --name tt-module --provider random
Using config file: /Users/tstraub/.tfx.hcl
List Module Versions for Organization: firefly
╭─────────┬────────╮
│ VERSION │ STATUS │
├─────────┼────────┤
│ 0.0.1   │ ok     │
│ 0.0.0   │ ok     │
╰─────────┴────────╯
```

## `tfx registry module version create`

Create a version for a module in the Private Registry.

This command will default to the working directory for the Terraform code.

Alternatively setting the `--directory` flag will upload that directory.

**Basic Example**

```sh
$ tfx registry module version create --name tt-module --provider random --version 0.0.1
Using config file: /Users/tstraub/.tfx.hcl
Create Module Version for Organization: firefly
Module Created, Uploading... 
Module Created: 
ID:      mod-3Gi7SGzjESjcUezn
Created: 2022-08-19T22:09:36.656Z
```

**Directory Example**

```sh
$ tfx registry module version create --name tt-module --provider random --version 0.0.2 --directory ./module/tt-module/
Using config file: /Users/tstraub/.tfx.hcl
Create Module Version for Organization: firefly
Module Created, Uploading... 
Module Created: 
ID:      mod-3Gi7SGzjESjcUezn
Created: 2022-08-19T22:11:57.765Z
```

## `tfx registry module version delete`

Delete a version of a module in the Private Registry.

**Example**

```sh
$ tfx registry module version delete --name tt-module --provider random -v 0.0.1
Using config file: /Users/tstraub/.tfx.hcl
Delete Module Version for Organization: firefly
Module Version Deleted: tt-module
Status: Success
```

## `tfx registry module version download`

Download a version of a module in the Private Registry.

**Temp Folder Example**

```sh
$ tfx registry module version download --name tt-module --provider random --version 0.0.1
Using config file: /Users/tstraub/.tfx.hcl
Downloading Module Version: tt-module
Directory not supplied, creating a temp directory 
Module Version Found, download started... 
Status:    Success
Directory: /var/folders/99/srh_6psj6g5520gwyv8v3nbw0000gn/T/slug2213227994/tt-module/random/0.0.1/
```

**Specific Folder Example**

```sh
$ tfx registry module version download --name tt-module --provider random --version 0.0.1 --directory ./local
Using config file: /Users/tstraub/.tfx.hcl
Downloading Module Version: tt-module
Directory not supplied, creating a temp directory 
Module Version Found, download started... 
Status:    Success
Directory: ./local/tt-module/random/0.0.1/
```
