# Development

VS Code Plugin -> https://marketplace.visualstudio.com/items?itemName=golang.Go

go install github.com/go-delve/delve/cmd/dlv@latest

https://github.com/golang/vscode-go/blob/master/docs/debugging.md#launch-configuration


```
go get -u github.com/spf13/cobra
go get github.com/spf13/cobra/cobra
mkdir tfx && cd tfx
~/go/bin/cobra init --pkg-name github.com/straubt1/tfx
~/go/bin/cobra add plan


export GOPATH=/usr/local/bin/go
```

## PMR

https://www.terraform.io/docs/registry/api.html


https://app.terraform.io/api/registry/v1/modules/terraform-tom/
https://firefly.tfe.rocks/.well-known/terraform.json
-> 	modules.v1	"/api/registry/v1/modules/"
https://firefly.tfe.rocks/api/registry/v1/modules/firefly/


https://firefly.tfe.rocks/api/v2/registry-modules/show/firefly/test/none

test2-ID:"mod-BAYTcvZTjMZ5dSzk"


https://firefly.tfe.rocks/v1/modules


Deleting a module via name on the PMR, will delete ALL modules with that name, regardless of provider