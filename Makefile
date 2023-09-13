# For local development only
build:
	go build -o bin/main main.go
	go build -v \
	-ldflags="-X '$(shell git remote get-url --push origin | sed 's/:/\//g' | sed 's/.git//g')/version.Version=9.9.9' \
	-X '$(shell git remote get-url --push origin | sed 's/:/\//g' | sed 's/.git//g')/version.Prerelease=alpha' \
	-X '$(shell git remote get-url --push origin | sed 's/:/\//g' | sed 's/.git//g')/version.Build=local' \
	-X '$(shell git remote get-url --push origin | sed 's/:/\//g' | sed 's/.git//g')/version.BuiltBy=$(shell git config --global  --get github.user)' \
	-X '$(shell git remote get-url --push origin | sed 's/:/\//g' | sed 's/.git//g')/version.Date=$(shell date)'"
	rm -rf bin/

update:
	go get -u
	go mod tidy

upgrade-go-mac:
	brew upgrade go

site-local:
	mkdocs serve -f site/mkdocs.yml

format:
	go fmt
