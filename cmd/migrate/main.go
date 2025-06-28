package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/andrey/epoch-server/tools/migration"
	"github.com/go-pkgz/lgr"
)

func main() {
	var (
		subgraphEndpoint = flag.String("subgraph", "http://localhost:8000/subgraphs/name/rewards", "Subgraph endpoint URL")
		vaultID          = flag.String("vault", "", "Vault ID to migrate (required)")
		snapshotBlock    = flag.String("block", "0", "Snapshot block number")
		dryRun           = flag.Bool("dry-run", true, "Run in dry-run mode (no actual updates)")
		verbose          = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *vaultID == "" {
		fmt.Fprintf(os.Stderr, "Error: vault ID is required\n")
		flag.Usage()
		os.Exit(1)
	}

	logger := lgr.NoOp
	if *verbose {
		logger = lgr.New(lgr.Msec, lgr.LevelBraces, lgr.CallerFile, lgr.CallerFunc).Logf
	}

	snapshotBlockNumber, ok := new(big.Int).SetString(*snapshotBlock, 10)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: invalid snapshot block number: %s\n", *snapshotBlock)
		os.Exit(1)
	}

	config := migration.MigrationConfig{
		SubgraphEndpoint:    *subgraphEndpoint,
		SnapshotBlockNumber: snapshotBlockNumber,
		VaultID:             *vaultID,
		DryRun:              *dryRun,
	}

	graphClient := graph.NewClient(*subgraphEndpoint)
	subsidizerClient := subsidizer.NewClient(logger)

	migrationService := migration.NewMigrationService(
		graphClient,
		subsidizerClient,
		logger,
		config,
	)

	ctx := context.Background()

	status, err := migrationService.GetMigrationStatus(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting migration status: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Migration Status:\n")
	for key, value := range status {
		fmt.Printf("  %s: %v\n", key, value)
	}
	fmt.Println()

	if *dryRun {
		fmt.Println("Running in DRY-RUN mode. No actual changes will be made.")
	}

	fmt.Printf("Starting subsidy initialization for vault %s...\n", *vaultID)

	err = migrationService.InitializeSubsidies(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration completed successfully!")
}
