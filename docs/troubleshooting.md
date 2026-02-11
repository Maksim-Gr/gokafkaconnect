# Troubleshooting

## "No Kafka Connect URL configured"

Run:

```bash
gk config configure
```

This creates `~/.gokafkaconnect/config.yaml`.

## "Failed to load config"

Your config file may be missing or invalid YAML. Re-run configuration:

```bash
gk config configure
```

## Connection errors (connection refused, timeout)

- Verify the Kafka Connect URL is correct and reachable.
- If you omitted the scheme, `gk` assumes `http://`.
- Check firewall or network policies between your machine and the cluster.

## Authentication failures (401/403)

- Confirm the username/password in `~/.gokafkaconnect/config.yaml`.
- If your cluster does not require auth, leave `username` and `password` empty.

## "No connectors found"

Your Kafka Connect cluster has no connectors, or your user has no access. Create one and retry.

## Delete errors with "rebalance is in process"

Kafka Connect reports a conflict (rebalance). Wait a bit and retry:

```bash
gk connector delete --connector my-connector
```

## Backup file not created

Ensure the backup output directory is writable:

```bash
gk config backup --dir ./backups
```
