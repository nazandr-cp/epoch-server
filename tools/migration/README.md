# Subsidy Migration Tool

This tool initializes `lastEffectiveValue` for existing `AccountSubsidiesPerCollection` records and publishes the first cumulative Merkle root for debt subsidization.

## Overview

The migration script performs the following operations:

1. **Fetches existing data**: Connects to the subgraph to retrieve all `AccountSubsidiesPerCollection` records for a specific vault.

2. **Computes `lastEffectiveValue`**: For each record, calculates `lastEffectiveValue = weightedBalance (if non-zero) + currentBorrowU` using the account's current borrow balance.

3. **Updates records**: Updates each entity's `lastEffectiveValue` in the subgraph or staging table (placeholder implementation).

4. **Generates Merkle root**: Creates the first cumulative Merkle root using `secondsAccumulated` data, builds a Merkle tree over `account → secondsAccumulated` mappings, and calls `DebtSubsidizer.UpdateMerkleRoot(vaultId, root)`.

5. **Ensures idempotence**: Can be run multiple times without duplicating work.

## Usage

### Command Line Interface

```bash
# Build the migration tool
go build ./cmd/migrate

# Run in dry-run mode (default, no actual changes)
./migrate -vault=0x1234567890abcdef -subgraph=http://localhost:8000/subgraphs/name/rewards

# Run with actual updates
./migrate -vault=0x1234567890abcdef -subgraph=http://localhost:8000/subgraphs/name/rewards -dry-run=false

# With verbose logging
./migrate -vault=0x1234567890abcdef -verbose

# Specify snapshot block
./migrate -vault=0x1234567890abcdef -block=12345678
```

### Command Line Options

- `-vault`: (Required) Vault ID to migrate
- `-subgraph`: Subgraph endpoint URL (default: `http://localhost:8000/subgraphs/name/rewards`)
- `-block`: Snapshot block number (default: `0`)
- `-dry-run`: Run in dry-run mode - no actual updates (default: `true`)
- `-verbose`: Enable verbose logging (default: `false`)

### Programmatic Usage

```go
package main

import (
    "context"
    "math/big"
    
    "github.com/andrey/epoch-server/internal/clients/graph"
    "github.com/andrey/epoch-server/internal/clients/subsidizer"
    "github.com/andrey/epoch-server/tools/migration"
    "github.com/go-pkgz/lgr"
)

func main() {
    config := migration.MigrationConfig{
        SubgraphEndpoint:    "http://localhost:8000/subgraphs/name/rewards",
        SnapshotBlockNumber: big.NewInt(12345678),
        VaultID:             "0x1234567890abcdef",
        DryRun:              false,
    }

    graphClient := graph.NewClient(config.SubgraphEndpoint)
    subsidizerClient := subsidizer.NewClient(lgr.NoOp)
    
    migrationService := migration.NewMigrationService(
        graphClient,
        subsidizerClient,
        lgr.NoOp,
        config,
    )

    ctx := context.Background()
    err := migrationService.InitializeSubsidies(ctx)
    if err != nil {
        // handle error
    }
}
```

## Implementation Details

### Data Structures

- **`AccountSubsidyRecord`**: Represents an account's subsidy state with fields for account address, weighted balance, current borrow amount, last effective value, and accumulated seconds.

- **`MerkleLeaf`**: Represents a leaf in the Merkle tree with account address and accumulated seconds.

- **`MigrationConfig`**: Configuration for the migration including subgraph endpoint, snapshot block, vault ID, and dry-run flag.

### Key Functions

- **`fetchAccountSubsidyRecords()`**: Queries the subgraph for existing records
- **`computeLastEffectiveValues()`**: Calculates the effective values using the formula
- **`updateSubgraphRecords()`**: Updates records in the subgraph (placeholder)
- **`generateMerkleRoot()`**: Builds and returns the Merkle tree root
- **`buildMerkleTree()`**: Recursive Merkle tree construction algorithm

### Merkle Tree Construction

The Merkle tree is built using a standard binary tree approach:

1. Sort accounts alphabetically for deterministic ordering
2. Create leaf nodes from `account:secondsAccumulated` pairs
3. Hash each leaf using SHA-256
4. Recursively combine hashes until reaching the root
5. Handle odd numbers of leaves by carrying the last leaf to the next level

## Testing

Run the test suite:

```bash
go test ./tools/migration -v
```

The tests cover:

- **Value computation**: Validates the `lastEffectiveValue` calculation logic
- **Merkle root generation**: Tests empty, single, and multiple record scenarios
- **Deterministic behavior**: Ensures consistent results regardless of input order
- **Tree construction**: Validates the Merkle tree building algorithm
- **End-to-end flow**: Tests the complete migration process with mocks

## Error Handling

The migration script includes comprehensive error handling for:

- Invalid vault IDs or block numbers
- Network connectivity issues with the subgraph
- Malformed data from GraphQL responses
- Failed Merkle root updates
- Arithmetic errors in value calculations

## Idempotence

The migration can be safely run multiple times. The `IsIdempotent()` method checks if the migration has already been completed for a given vault, preventing duplicate work.

## Security Considerations

- All GraphQL queries use parameterized inputs to prevent injection attacks
- Merkle tree construction uses cryptographically secure SHA-256 hashing
- The dry-run mode allows safe testing before actual execution
- Comprehensive input validation prevents malformed data processing