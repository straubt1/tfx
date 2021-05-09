# TFx CLI

CLI to wrap the Terraform Cloud and Terraform Enterprise API for common tasks.
The primary focus initially is the API-Driven workflow.

## Run Workflow

**Create a Plan**

```sh
# Create a speculative plan that can not be applied
tfx plan -w tfx-test -s

# Create a plan that can be applied
tfx plan -w tfx-test
```

**Create an Apply**

```sh
tfx apply -r <run-id>
```

### Example

```sh
$tfx plan -w tfx-test
Using config file: /Users/tstraub/tfx/.tfx.hcl
Creating new Config Version
Workspace Run Created, Run Id: run-pEWkbDy7aNqBztNQ Config Version: cv-J7Apwj3fYojsBeNu
Navigate: https://firefly.tfe.rocks/app/firefly/workspaces/tfx-test/runs/run-pEWkbDy7aNqBztNQ

------------------------------------------------------------------------
Terraform v0.14.11
Configuring remote state backend...
Initializing Terraform configuration...

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # random_pet.one will be created
  + resource "random_pet" "one" {
      + id        = (known after apply)
      + length    = 2
      + separator = "-"
    }

  # random_pet.two will be created
  + resource "random_pet" "two" {
      + id        = (known after apply)
      + length    = 2
      + separator = "-"
    }

Plan: 2 to add, 0 to change, 0 to destroy.
------------------------------------------------------------------------
Run Complete: run-pEWkbDy7aNqBztNQ
```

```sh
$ tfx apply -r run-pEWkbDy7aNqBztNQ                                                                                                                      

Using config file: /Users/tstraub/tfx/.tfx.hcl
Workspace Apply Created, Apply Id: apply-38E9GEC14FeAJcZv
Navigate: https://firefly.tfe.rocks/app/firefly/workspaces//runs/run-pEWkbDy7aNqBztNQ

Terraform v0.14.11
random_pet.one: Creating...
random_pet.two: Creating...
random_pet.one: Creation complete after 0s [id=rational-stud]
random_pet.two: Creation complete after 0s [id=vast-silkworm]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
Apply Complete: apply-38E9GEC14FeAJcZv
```

## Future Commands

- `tfx run list`, List all runs for a workspace
- `tfx cv list`, List all configuration versions for a workspace
- `tfe tfversions`
  - `list`, list all Terraform versions in TFE
  - `disable`, disable a Terraform version, -a flag to disable all
  - `enable`, enable a Terraform version
  - `create`, create a new Terraform version, upsert?
- `tfx pmr`
  - `list`, list all modules in the PMR
  - `create`, create a module in the PMR
  - `publish`, create a version of a module in the PMR, pushes code
  - `delete`, deletes a module/version from the PMR
- `tfe sentinel`
  - `list`, list policy sets
  - `create`, create a policy set
  - `delete`, deletes a policy set
  - `assign`, assigns a WS(s) to the policy set



## Feature Thoughts

- `tfx init` could still be valuable, maybe pull state file locally or verify that a Workspace exists and is ready (not locked)



## References

https://github.com/hashicorp/go-tfe

https://github.com/spf13/cobra#installing