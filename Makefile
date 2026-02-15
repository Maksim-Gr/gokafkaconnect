.PHONY: lint fmt test staticcheck test-coverage all

staticcheck:
	staticcheck -tests ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

all: fmt lint test staticcheck