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

TODO: Add setup and running instructions