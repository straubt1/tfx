# For local development only
build:
	go build -o bin/main main.go
	go build -v \
	-ldflags="-X 'github.com/straubt1/tfx/version.Version=9.9.9' \
	-X 'github.com/straubt1/tfx/version.Prerelease=alpha' \
	-X 'github.com/straubt1/tfx/version.Build=local' \
	-X 'github.com/straubt1/tfx/version.BuiltBy=tstraub' \
	-X 'github.com/straubt1/tfx/version.Date=$(shell date)'"
	rm -rf bin/

update:
	go get -u
	go mod tidy

site-local:
	mkdocs serve -f site/mkdocs.yml

format:
	go fmt
