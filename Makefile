.PHONY: test

setup: dep
	mkdir -p data

dep:
	which dep || go get -u github.com/golang/dep/cmd/dep
	dep ensure
	dep ensure -update

generate:
	find ./data -type f | xargs cat | go run cmd/runner.go > result.csv

test:
	go test ./...

run:
	$(MAKE) -C docker $@
run-test:
	$(MAKE) -C docker $@
