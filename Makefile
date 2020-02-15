.PHONY: test

setup: tools mod generate
	mkdir -p data

tools:
	which go-assets-builder || go get -u github.com/jessevdk/go-assets-builder

mod:
	go mod download

generate:
	go generate ./...

anki:
	find ./data -type f | xargs cat | go run cmd/cli/main.go --dictionary=wiktionary > result.csv

test:
	go test ./...

heroku: setup
	go build ./cmd/server/main.go

run:
	$(MAKE) -C docker $@

run-test:
	$(MAKE) -C docker $@
