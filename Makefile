# For local development only
build:
	BUILD_SHA="local" BUILD_DATE=$(date) BUILT_BY="me" \
	goreleaser release --snapshot --clean

update:
	go get -u
	go mod tidy

upgrade-go-mac:
	brew upgrade go

site-local:
	mkdocs serve -f site/mkdocs.yml

format:
	go fmt
