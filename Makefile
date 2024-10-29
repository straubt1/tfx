# For local development only
build:
	go build -o bin/main main.go
	go build -v \
	-ldflags="-X '$(shell git remote get-url --push origin | sed 's/https\?:\/\///' | sed 's/\.git//g')/version.Version=0.1.4' \
	-X '$(shell git remote get-url --push origin | sed 's/https\?:\/\///' | sed 's/\.git//g')/version.Prerelease=alpha' \
	-X '$(shell git remote get-url --push origin | sed 's/https\?:\/\///' | sed 's/\.git//g')/version.Build=local' \
	-X '$(shell git remote get-url --push origin | sed 's/https\?:\/\///' | sed 's/\.git//g')/version.BuiltBy=$(shell git config --global  --get github.user)' \
	-X '$(shell git remote get-url --push origin | sed 's/https\?:\/\///' | sed 's/\.git//g')/version.Date=$(shell date)'"
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
