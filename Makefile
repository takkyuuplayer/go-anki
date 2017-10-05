.PHONY: test

setup:
	mkdir -p data

generate:
	find ./data -type f | xargs cat | go run runner.go > result.csv

test:
	go test ./...

run:
	make -C docker $@
run-test:
	make -C docker $@
