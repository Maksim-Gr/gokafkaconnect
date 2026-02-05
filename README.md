# kc CLI

A command-line interface for managing Kafka Connect connectors via the Kafka Connect REST API.  
`kc` focuses on providing a fast, simple, and interactive CLI experience for day-to-day connector operations.

---

## Overview

`kc` is a Go-based CLI tool designed to interact with Kafka Connect clusters.  
It creates a lightweight client for the Kafka Connect REST API and exposes common connector management operations through an intuitive command-line interface.

The tool is intended for developers and operators who want a straightforward way to list, inspect, back up, create, and delete connectors without manually interacting with REST endpoints.

---

## Features

- List running Kafka Connect connectors
- View connector configurations
- Create connectors from predefined templates
- Delete existing connectors
- Back up connector configurations to JSON files
- Interactive CLI prompts (arrow-key navigation)
- Simple configuration-driven setup

---

## Configuration

`kc` requires a configuration file to locate and connect to a Kafka Connect cluster.

The configuration file defines at minimum:
- Kafka Connect REST API URL

The config is loaded at runtime using the internal configuration loader.  
(Example configuration and documentation will be expanded as the project evolves.)

---

## Installation

Clone the repository and build the binary locally:

```bash
git clone https://github.com/Maksim-Gr/gokafkaconnect.git
cd gokafkaconnect
go build -o kc
```

Run the CLI:

```bash
./kc
```

---

## Usage

The CLI exposes multiple commands to interact with Kafka Connect.

Typical workflows include:
- Listing available connectors
- Selecting a connector interactively
- Viewing or backing up connector configuration
- Creating or deleting connectors

Run the following to explore available commands:

```bash
./kc --help
```

---

## Backup Example

The `backup` command retrieves all connector configurations from the Kafka Connect cluster and stores them in a timestamped JSON file:

```bash
./kc backup --dir ./backup
```

This allows connector configurations to be versioned, reviewed, or restored later.

---

## Planned Improvements

- Update existing connector configurations
- Load connector definitions dynamically from YAML / JSON
- Improved output formatting and status reporting
- Enhanced error handling and diagnostics
- Expanded configuration options (authentication, TLS, etc.)

---

## Project Status

This project is under active development and is currently focused on core connector lifecycle operations.  
Breaking changes may occur while APIs and internal structure are refined.

---

## Contributing & Feedback

`kc` is a personal project created to simplify connector management for my own use.  
Bug reports, feedback, and contributions are welcome.

If you encounter issues or have suggestions:
- Open a GitHub issue
- Fork the project and submit a pull request

---

## References

- Kafka Connect REST API documentation:  
  https://docs.confluent.io/platform/current/connect/references/restapi.html
