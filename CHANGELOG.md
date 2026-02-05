# CHANGELOG

---

## Unreleased

Initial development of **kc**.

### Added
- Kafka Connect REST API client abstraction
- Connector operations:
    - List connectors
    - View connector configuration (raw and JSON)
    - Create connectors from predefined templates
    - Delete connectors
    - Backup connector configurations to timestamped JSON files
- Task operations:
    - List tasks for a connector
    - Get task status
    - Restart a task
- Config operations:
    - Configure Kafka Connect URL
    - Show current configuration
- Interactive CLI prompts for connector/task selection
- Configuration-driven Kafka Connect URL loading

### Changed
- CLI commands reorganized into subdirectories/packages (`cmd/config`, `cmd/connector`, `cmd/task`) for clearer separation

### Fixed
- Configuration file resolution to avoid failures when running from different working directories / build contexts

### Breaking Changes
- Command layout changed due to CLI package reorganization (subcommands moved under `config`, `connector`, `task`)

---

_This project is under active development. Versions and release notes will be added once the first stable release is published._