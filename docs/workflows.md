# Common workflows

## Initial setup and verify

```bash
gk config configure

gk connector list
```

## Create a connector from a JSON file

```bash
gk connector create --file ./connector.json
```

## Create a connector interactively

```bash
gk connector create
```

Follow the prompts to fill required fields and optionally submit the config.

## Inspect a connector config

```bash
gk connector list
```

Pick a connector when prompted to print its config.

## Check connector health

```bash
gk connector health-check
```

## List tasks and check a task status

```bash
gk task list --connector my-connector

gk task get --connector my-connector --id 0
```

## Restart a task

```bash
gk task restart --connector my-connector --id 0
```

## Backup all connector configs

```bash
gk config backup --dir ./backups
```

## Dry run a change

```bash
gk --dry-run task restart --connector my-connector --id 0
```
