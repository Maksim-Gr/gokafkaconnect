# Getting started

## Prerequisites

- Go installed (if building from source)
- Access to a Kafka Connect cluster (its REST API URL)

## Build from source

From the repository root:

```bash
go build -o gk .
```

## Run

```bash
./gk --help
```

## First-time setup

`gk` needs a Kafka Connect URL (and optional basic auth). You can configure it explicitly:

```bash
./gk config configure
```

If no configuration exists, `gk` will prompt you automatically on the next run.

## Quick verification

```bash
./gk connector list
```

If you get an error, jump to [Troubleshooting](./troubleshooting.md).
