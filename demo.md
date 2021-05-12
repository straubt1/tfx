# Demo



```sh
$ tfx pmr list
Using config file: /Users/tstraub/tfx/.tfx.hcl
╭──────┬──────────┬────┬───────────╮
│ NAME │ PROVIDER │ ID │ PUBLISHED │
├──────┼──────────┼────┼───────────┤
╰──────┴──────────┴────┴───────────╯

$ go run . pmr create --name my-module --provider aws 
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