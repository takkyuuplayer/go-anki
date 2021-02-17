.PHONY: test

setup: mod generate

mod:
	go mod download

generate:
	go generate ./...

test:
	go clean -testcache
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...

lint: golint gocyclo

golint:
	which golint || go get -u -v golang.org/x/lint/golint
	go list ./... | xargs golint

gocyclo:
	which gocyclo || go get -u -v github.com/fzipp/gocyclo/cmd/gocyclo
	gocyclo -over 20 .

heroku: setup
	go build ./cmd/server/main.go
