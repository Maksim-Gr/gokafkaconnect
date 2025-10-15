.PHONY: lint fmt test staticcheck all

staticcheck:
	staticcheck -tests ./...

lint:
	golangci-lint run
fmt:
	gofmt -s -w .
test:
	go test -v ./...

all: fmt lint test staticcheck