# Configuration

`gk` needs to know how to reach your Kafka Connect REST API and (optionally) basic auth credentials.

## Where configuration lives

`gk` stores its config at:

- `~/.gokafkaconnect/config.yaml`

You can create or update it via the interactive command:

```bash
./gk config configure
```

## Config format

`config.yaml` is YAML and follows this structure:

```yaml
kafkaConnect:
  url: http://localhost:8083
  username: ""
  password: ""
```

Notes:
- `username`/`password` are optional. Leave them empty for no auth.
- If you enter a URL without a scheme, `gk` assumes `http://`.

## View current config

```bash
./gk config show-config
```

## Dry run

Some commands support `--dry-run` (global flag) to show what would happen without making changes:

```bash
./gk --dry-run config configure
```
