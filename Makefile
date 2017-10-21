.PHONY: test

setup:
	mkdir -p data

generate:
	find ./data -type f | xargs cat | go run cmd/runner.go > result.csv

test:
	go test ./...

run:
	@cd docker && $(MAKE) run
run-test:
	@cd docker && $(MAKE) run-test
