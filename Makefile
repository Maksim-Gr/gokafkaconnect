.PHONY: lint fmt test

lint:
	golangci-lint run
fmt:
	gofmt -s -w .
test:
	go test ./..

all: fmt lint test