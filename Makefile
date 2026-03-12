.PHONY: lint fmt test staticcheck test-coverage all dev-up dev-down dev-logs dev-status

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

dev-up:
	docker compose up -d

dev-down:
	docker compose down -v

dev-logs:
	docker compose logs -f kafka-connect

dev-status:
	curl -s http://localhost:8083/connectors | jq .