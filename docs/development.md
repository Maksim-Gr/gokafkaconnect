# Development

## Build

```bash
go build -o gk .
```

## Run locally

```bash
./gk --help
```

## Tests

```bash
go test ./...
```

## Lint and static checks

```bash
golangci-lint run
staticcheck -tests ./...
```

## Format

```bash
gofmt -s -w .
```
