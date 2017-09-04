RESULT=$(shell date '+%s')

.PHONY: test

setup:
	mkdir -p data result

generate:
	find ./data -type f | xargs cat | go run runner.go > result/${RESULT}.tsv

test:
	go test ./...

run:
	@cd docker && $(MAKE) run
run-test:
	@cd docker && $(MAKE) run-test
