# Epoch Server

A Go server for managing epoch-based subsidy distribution in the lend.fam ecosystem.

## Overview

This server handles:
- Epoch management and transitions
- Subsidy calculation and distribution
- Integration with Ethereum smart contracts
- GraphQL subgraph data processing
- Scheduled epoch operations

## Project Structure

```
epoch-server/
├── cmd/server/           # Application entrypoint
├── internal/
│   ├── config/          # Configuration loading
│   ├── log/             # LGR logger setup
│   ├── clients/         # External service clients
│   │   ├── graph/       # Subgraph GraphQL client
│   │   └── contract/    # Ethereum contract client
│   ├── handlers/        # HTTP handlers using Chi
│   ├── service/         # Epoch management & subsidy logic
│   └── scheduler/       # Scheduled epoch runner
├── pkg/contracts/       # Go bindings for smart contracts
├── configs/             # Configuration files
└── scripts/             # Helper scripts
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Access to an Ethereum RPC endpoint (ApeChain Curtis testnet or local)
- Running subgraph instance (optional for full functionality)

### Installation

1. Navigate to the epoch-server directory:
```bash
cd epoch-server
```

2. Install dependencies:
```bash
go mod download
```

3. Build the server:
```bash
go build ./cmd/server
```

## Configuration

The server uses YAML configuration files. The default configuration is located at `configs/config.yaml`.

### Environment Variables

Set the following environment variables:

```bash
# Required: Ethereum private key for transaction signing
export ETHEREUM_PRIVATE_KEY="your_private_key_here"

# Optional: Override config file path
export CONFIG_PATH="configs/config.yaml"
```

### Configuration Options

Key configuration sections in `config.yaml`:

- **server**: HTTP server settings (host, port)
- **ethereum**: Blockchain connection (RPC URL, gas settings)
- **subgraph**: GraphQL endpoint for querying indexed data
- **scheduler**: Automated epoch processing schedule
- **contracts**: Smart contract addresses
- **features**: Development flags (dry_run, debug_mode)

## Usage

### Running the Server

Start the epoch server:

```bash
# Using default config
./server

# Using custom config
./server -config configs/config.yaml
```

### Database Migrations

Run database migrations (if using persistent storage):

```bash
go run ./cmd/migrate
```

### API Endpoints

The server exposes REST API endpoints:

- `GET /health` - Health check
- `GET /epoch/current` - Get current epoch information
- `POST /epoch/process` - Manually trigger epoch processing
- `GET /subsidies/{address}` - Get subsidy eligibility for an address

### Development Mode

For development and testing:

1. Enable dry run mode in config:
```yaml
features:
  dry_run: true
```

2. Use debug logging:
```yaml
logging:
  level: "debug"
```

### Scheduled Operations

The server automatically processes epochs based on the scheduler configuration:

```yaml
scheduler:
  interval: "1h"  # Process every hour
  enabled: true
  timezone: "UTC"
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Smart Contract Integration

After smart contract updates, regenerate Go bindings:

```bash
# From project root
./scripts/generate_bindings.sh
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Ensure RPC endpoint is accessible
2. **Transaction failed**: Check gas settings and account balance
3. **Subgraph errors**: Verify subgraph is deployed and synced
4. **Permission denied**: Ensure private key has necessary permissions

### Logs

Check server logs for detailed error information. Logs are output to stdout in JSON format by default.