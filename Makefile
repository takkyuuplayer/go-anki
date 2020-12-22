.PHONY: test

setup: tools mod generate

tools:
	which statik || go get github.com/rakyll/statik

mod:
	go mod download

generate:
	go generate ./...

test:
	go clean -testcache
	go test -race -covermode=atomic -coverprofile=coverage.txt ./...

lint: gocyclo

gocyclo:
	which gocyclo || go get -u -v github.com/fzipp/gocyclo/cmd/gocyclo
	gocyclo -over 20 .

heroku: setup
	go build ./cmd/server/main.go
