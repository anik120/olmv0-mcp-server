# OLM v0 MCP Server

A Model Context Protocol (MCP) server for Operator Lifecycle Manager (OLM) v0, providing AI assistants with the ability to interact with OLM resources in Kubernetes clusters.

## Overview

The OLM v0 MCP Server exposes OLM-specific resources and operations through the Model Context Protocol, allowing AI assistants to help users manage operators in their Kubernetes clusters. This server follows the same pattern as the [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server) but focuses specifically on OLM v0 resources.

## Features

- **ClusterServiceVersion (CSV) Operations**: List and inspect operator deployments
- **Subscription Management**: View operator subscriptions and update channels
- **CatalogSource Inspection**: Check catalog sources and their health
- **InstallPlan Analysis**: Review planned operator installations and updates
- **Read-only Mode**: Safe operation with no write capabilities by default
- **Configurable Toolsets**: Enable/disable specific functionality groups
- **Multi-cluster Support**: Connect to any Kubernetes cluster with OLM

## Available Tools

### ClusterServiceVersion Tools
- `list_csvs`: List ClusterServiceVersions in a namespace
- `get_csv`: Get detailed information about a specific ClusterServiceVersion

### Subscription Tools
- `list_subscriptions`: List Subscriptions in a namespace
- `get_subscription`: Get detailed information about a specific Subscription

### CatalogSource Tools
- `list_catalog_sources`: List CatalogSources in a namespace
- `get_catalog_source`: Get detailed information about a specific CatalogSource

### InstallPlan Tools
- `list_install_plans`: List InstallPlans in a namespace
- `get_install_plan`: Get detailed information about a specific InstallPlan

### General Tools
- `list_tools`: Show available tools and their parameters

## Installation

### From Source

```bash
git clone https://github.com/operator-framework/operator-lifecycle-manager.git
cd operator-lifecycle-manager/olmv0-mcp-server
make build
```

### Using Docker

```bash
docker build -t olmv0-mcp-server .
```

## Usage

### Basic Usage

```bash
# Start the server with default settings
./bin/olmv0-mcp-server

# Start on a specific port
./bin/olmv0-mcp-server --port 8080

# Use a specific kubeconfig
./bin/olmv0-mcp-server --kubeconfig /path/to/kubeconfig

# Enable only specific toolsets
./bin/olmv0-mcp-server --toolsets csv,subscription
```

### Command Line Options

- `--port, -p`: HTTP/SSE server port (default: 8080)
- `--kubeconfig`: Path to kubeconfig file (default: $HOME/.kube/config)
- `--read-only`: Prevent write operations (default: true)
- `--toolsets`: Enable specific toolsets (default: csv,subscription,catalog,installplan)

### Docker Usage

```bash
# Run with default settings
docker run --rm -p 8080:8080 \
  -v ~/.kube/config:/root/.kube/config:ro \
  olmv0-mcp-server

# Run with custom kubeconfig
docker run --rm -p 8080:8080 \
  -v /path/to/kubeconfig:/kubeconfig:ro \
  olmv0-mcp-server --kubeconfig /kubeconfig
```

## API Examples

### List ClusterServiceVersions

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "method": "list_csvs",
    "params": {
      "namespace": "operators"
    }
  }'
```

### Get Specific Subscription

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "method": "get_subscription",
    "params": {
      "name": "my-operator",
      "namespace": "operators"
    }
  }'
```

### List Available Tools

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "method": "list_tools",
    "params": {}
  }'
```

## Development

### Prerequisites

- Go 1.24.4 or later
- Access to a Kubernetes cluster with OLM installed
- kubectl configured for cluster access

### Building

```bash
# Download dependencies
make deps

# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Project Structure

```
olmv0-mcp-server/
├── cmd/olmv0-mcp-server/     # Main application entry point
├── pkg/
│   ├── client/               # OLM Kubernetes client wrappers
│   ├── server/               # HTTP server and MCP request handling
│   ├── tools/                # Tool implementations for each resource type
│   └── types/                # Common types and interfaces
├── docs/                     # Documentation
├── Dockerfile               # Container build file
├── Makefile                # Build automation
└── README.md               # This file
```

## Configuration

The server uses the standard Kubernetes client configuration:

1. **In-cluster**: If running inside a Kubernetes cluster, uses the service account token
2. **Kubeconfig**: Uses the kubeconfig file specified by `--kubeconfig` flag
3. **Default**: Falls back to `$HOME/.kube/config`

## Security Considerations

- **Read-only by default**: The server operates in read-only mode by default
- **No write operations**: Currently only supports read operations (list, get, inspect)
- **Cluster access**: Requires valid Kubernetes credentials with appropriate RBAC permissions
- **No authentication**: The HTTP server does not implement authentication (intended for local use)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

## Related Projects

- [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server) - General Kubernetes MCP server
- [Operator Lifecycle Manager](https://github.com/operator-framework/operator-lifecycle-manager) - Main OLM project
- [Model Context Protocol](https://spec.modelcontextprotocol.io/) - MCP specification