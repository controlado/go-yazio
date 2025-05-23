.PHONY: run test

run:
	go run ./cmd/...

test:
	@echo Running tests
	go test -coverpkg=./... -covermode=atomic ./...