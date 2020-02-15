.PHONY: test

setup: tools mod generate
	mkdir -p data

tools:
	@which statik || go get github.com/rakyll/statik

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
