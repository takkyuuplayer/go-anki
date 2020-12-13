.PHONY: test

setup: tools mod generate
	mkdir -p data

tools:
	which statik || go get github.com/rakyll/statik

mod:
	go mod download

generate:
	go generate ./...

anki:
	find ./data -type f | xargs cat | go run cmd/cli/main.go --dictionary=wiktionary > result.csv

test:
	go clean -testcache
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...

lint: golint gocyclo

golint:
	which golint || go get -u -v golang.org/x/lint/golint
	go list ./... | xargs golint

gocyclo:
	which gocyclo || go get -u -v github.com/fzipp/gocyclo
	find . -maxdepth 1 -mindepth 1 -type d -regex "\.\/[a-z].*" | grep -v vendor | xargs gocyclo -over 15

heroku: setup
	go build ./cmd/server/main.go

run:
	$(MAKE) -C docker $@

run-test:
	$(MAKE) -C docker $@
