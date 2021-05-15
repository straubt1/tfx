# Demo

## Plan/Apply

```sh
$ tfx plan -w tfx-test-tfc -d ./terraform
Using config file: /Users/tstraub/tfx/.tfx.hcl
Remote Terraform Plan, speculative plan: false  destroy plan: false
Reading Workspace tfx-test-tfc for ID... Found: ws-sr6nbVudgwchkFYf
Creating new Config Version ... ID: cv-KdJq8HyH49QGvpD9
Workspace Run Created, Run Id: run-x2qNtxb7aQELaCrJ Config Version: cv-KdJq8HyH49QGvpD9
Navigate: https://app.terraform.io/app/terraform-tom/workspaces/tfx-test-tfc/runs/run-x2qNtxb7aQELaCrJ

------------------------------------------------------------------------
Terraform v0.15.3
on linux_amd64
Configuring remote state backend...
Initializing Terraform configuration...

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
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
Cost estimation:
Waiting for cost estimate to complete...
Cost estimation:
Resources: 0 of 0 estimated
           $0.0/mo +$0.0
------------------------------------------------------------------------
Organization policy check:
================ Results for policy set: <empty policy set name> ===============

Sentinel Result: true

This result means that all Sentinel policies passed and the protected
behavior is allowed.

1 policies evaluated.

## Policy 1: this-policy-is-always-true (hard-mandatory)

Result: true

./this-policy-is-always-true.sentinel:1:1 - Rule "main"
  Value:
    true


Run Complete: run-x2qNtxb7aQELaCrJ

$ tfx apply -r run-x2qNtxb7aQELaCrJ
Using config file: /Users/tstraub/tfx/.tfx.hcl
Workspace Apply Created, Apply Id: apply-rF1wqfUAiFbwsvuf
Navigate: https://app.terraform.io/app/terraform-tom/workspaces//runs/run-x2qNtxb7aQELaCrJ

------------------------------------------------------------------------
Terraform v0.15.3
on linux_amd64
random_pet.two: Creating...
random_pet.one: Creating...
random_pet.one: Creation complete after 0s [id=up-guppy]
random_pet.two: Creation complete after 0s [id=heroic-vervet]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Apply Complete: apply-rF1wqfUAiFbwsvuf


$ tfx plan -w tfx-test-tfc -d ./terraform --destroy                                                                                                  
Using config file: /Users/tstraub/tfx/.tfx.hcl
Remote Terraform Plan, speculative plan: false  destroy plan: true
Reading Workspace tfx-test-tfc for ID... Found: ws-sr6nbVudgwchkFYf
Creating new Config Version ... ID: cv-4NyNQmmH3nghQGMG
Workspace Run Created, Run Id: run-ZDAEgRijY5KDkm1s Config Version: cv-4NyNQmmH3nghQGMG
Navigate: https://app.terraform.io/app/terraform-tom/workspaces/tfx-test-tfc/runs/run-ZDAEgRijY5KDkm1s

------------------------------------------------------------------------
Terraform v0.15.3
on linux_amd64
Configuring remote state backend...
Initializing Terraform configuration...
random_pet.one: Refreshing state... [id=up-guppy]
random_pet.two: Refreshing state... [id=heroic-vervet]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # random_pet.one will be destroyed
  - resource "random_pet" "one" {
      - id        = "up-guppy" -> null
      - length    = 2 -> null
      - separator = "-" -> null
    }

  # random_pet.two will be destroyed
  - resource "random_pet" "two" {
      - id        = "heroic-vervet" -> null
      - length    = 2 -> null
      - separator = "-" -> null
    }

Plan: 0 to add, 0 to change, 2 to destroy.

------------------------------------------------------------------------
Cost estimation:
Waiting for cost estimate to complete...
Cost estimation:
Resources: 0 of 0 estimated
           $0.0/mo +$0.0
------------------------------------------------------------------------
Organization policy check:
================ Results for policy set: <empty policy set name> ===============

Sentinel Result: true

This result means that all Sentinel policies passed and the protected
behavior is allowed.

1 policies evaluated.

## Policy 1: this-policy-is-always-true (hard-mandatory)

Result: true

./this-policy-is-always-true.sentinel:1:1 - Rule "main"
  Value:
    true


Run Complete: run-ZDAEgRijY5KDkm1s

$ tfx apply -r run-ZDAEgRijY5KDkm1s                                                                                                                   
Using config file: /Users/tstraub/tfx/.tfx.hcl
Workspace Apply Created, Apply Id: apply-czK687MwfkwsAhYJ
Navigate: https://app.terraform.io/app/terraform-tom/workspaces//runs/run-ZDAEgRijY5KDkm1s

------------------------------------------------------------------------
Terraform v0.15.3
on linux_amd64
random_pet.one: Destroying... [id=up-guppy]
random_pet.two: Destroying... [id=heroic-vervet]
random_pet.one: Destruction complete after 0s
random_pet.two: Destruction complete after 0s

Apply complete! Resources: 0 added, 0 changed, 2 destroyed.

Apply Complete: apply-czK687MwfkwsAhYJ
```

## PMR

```sh
$ tfx pmr list
Using config file: /Users/tstraub/tfx/.tfx.hcl
╭──────┬──────────┬────┬───────────╮
│ NAME │ PROVIDER │ ID │ PUBLISHED │
├──────┼──────────┼────┼───────────┤
╰──────┴──────────┴────┴───────────╯

$ tfx pmr create --name my-module --provider aws 
Using config file: /Users/tstraub/tfx/.tfx.hcl
Creating Module my-module/aws ...  Created with ID:  mod-fKErTKrJX2eZnESG

$ tfx pmr create version --name my-module --provider aws --moduleVersion 0.0.1
Using config file: /Users/tstraub/tfx/.tfx.hcl
Creating Module Version my-module/aws:0.0.1 ...  Uploading ...  Module Version Created

$ tfx pmr create version --name my-module --provider aws --moduleVersion 0.0.2
Using config file: /Users/tstraub/tfx/.tfx.hcl
Creating Module Version my-module/aws:0.0.2 ...  Uploading ...  Module Version Created

$ tfx pmr list
Using config file: /Users/tstraub/tfx/.tfx.hcl
╭───────────┬──────────┬─────────────────────────────┬──────────────────────────────────────╮
│ NAME      │ PROVIDER │ ID                          │ PUBLISHED                            │
├───────────┼──────────┼─────────────────────────────┼──────────────────────────────────────┤
│ my-module │ aws      │ firefly/my-module/aws/0.0.2 │ 2021-05-12 00:32:17.502172 +0000 UTC │
╰───────────┴──────────┴─────────────────────────────┴──────────────────────────────────────╯

$ tfx pmr show --name my-module --provider aws
Using config file: /Users/tstraub/tfx/.tfx.hcl
Showing Module my-module/aws... Found
ID:         mod-fKErTKrJX2eZnESG
Status:     setup_complete
Versions:   2
Created:    2021-05-12T00:29:35.282Z
Updated:    2021-05-12T00:30:08.959Z

$ tfx pmr show versions --name my-module --provider aws
Using config file: /Users/tstraub/tfx/.tfx.hcl
Showing Module my-module/aws... Found
╭─────────┬────────╮
│ VERSION │ STATUS │
├─────────┼────────┤
│ 0.0.2   │ ok     │
│ 0.0.1   │ ok     │
╰─────────┴────────╯

$ tfx pmr download --name my-module --provider aws --moduleVersion 0.0.1
Using config file: /Users/tstraub/tfx/.tfx.hcl
Downloading Module Version my-module/aws:0.0.1... Downloaded to Temp:  /var/folders/mk/l44pbn5x4bq2qv11vrbj755w0000gp/T/slug362667336

$ tfx pmr delete version --name my-module --provider aws --moduleVersion 0.0.1
Using config file: /Users/tstraub/tfx/.tfx.hcl
Deleting Module Version for my-module/aws:0.0.1... Deleted

$ tfx pmr show versions --name my-module --provider aws
Using config file: /Users/tstraub/tfx/.tfx.hcl
Showing Module my-module/aws... Found
╭─────────┬────────╮
│ VERSION │ STATUS │
├─────────┼────────┤
│ 0.0.2   │ ok     │
╰─────────┴────────╯            

$ tfx pmr delete --name my-module --provider aws
Using config file: /Users/tstraub/tfx/.tfx.hcl
Deleting Module for my-module... Deleted

$ tfx pmr list
Using config file: /Users/tstraub/tfx/.tfx.hcl
╭──────┬──────────┬────┬───────────╮
│ NAME │ PROVIDER │ ID │ PUBLISHED │
├──────┼──────────┼────┼───────────┤
╰──────┴──────────┴────┴───────────╯
```