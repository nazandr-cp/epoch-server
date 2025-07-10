package scheduler

import (
	"context"
	"time"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/go-pkgz/lgr"
)

//go:generate moq -out scheduler_mocks.go . EpochService SubsidyService

// EpochService interface for epoch operations
type EpochService interface {
	StartEpoch(ctx context.Context) error
}

// SubsidyService interface for subsidy operations
type SubsidyService interface {
	DistributeSubsidies(ctx context.Context, vaultId string) error
}

// Scheduler manages automated epoch operations
type Scheduler struct {
	epochService   EpochService
	subsidyService SubsidyService
	logger         lgr.L
	interval       time.Duration
	config         *config.Config
}
