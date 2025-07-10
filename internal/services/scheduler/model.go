package scheduler

import (
	"context"
	"time"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
)

//go:generate moq -out scheduler_mocks.go . EpochService SubsidyService

// EpochService interface for epoch operations
type EpochService interface {
	StartEpoch(ctx context.Context) (*epoch.StartEpochResponse, error)
}

// SubsidyService interface for subsidy operations
type SubsidyService interface {
	DistributeSubsidies(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error)
}

// Scheduler manages automated epoch operations
type Scheduler struct {
	epochService   epoch.Service
	subsidyService subsidy.Service
	logger         lgr.L
	interval       time.Duration
	config         *config.Config
}
