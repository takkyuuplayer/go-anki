.PHONY: test

setup: tools dep generate
	mkdir -p data

tools:
	which dep || go get -u github.com/golang/dep/cmd/dep
	which go-assets-builder || go get -u github.com/jessevdk/go-assets-builder

dep:
	dep ensure
	dep ensure -update

generate:
	go generate ./...

anki:
	find ./data -type f | xargs cat | go run cmd/cli/main.go --dictionary=wiktionary > result.csv

test:
	go test ./...

heroku: tools dep generate
	go build ./cmd/server/main.go

run:
	$(MAKE) -C docker $@

run-test:
	$(MAKE) -C docker $@
