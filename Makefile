.PHONY: test

test:
	go test ./...

run:
	@cd docker && $(MAKE) run
run-test:
	@cd docker && $(MAKE) run-test
