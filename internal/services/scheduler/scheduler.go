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

type EpochService interface {
	StartEpoch(ctx context.Context) error
}

type SubsidyService interface {
	DistributeSubsidies(ctx context.Context, vaultId string) error
}

type Scheduler struct {
	epochService   EpochService
	subsidyService SubsidyService
	logger         lgr.L
	interval       time.Duration
	config         *config.Config
}

func NewScheduler(epochService epoch.Service, subsidyService subsidy.Service, interval time.Duration, logger lgr.L, cfg *config.Config) *Scheduler {
	return &Scheduler{
		epochService:   epochService,
		subsidyService: subsidyService,
		logger:         logger,
		interval:       interval,
		config:         cfg,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Logf("INFO scheduler started with interval %v", s.interval)

	for {
		select {
		case <-ctx.Done():
			s.logger.Logf("INFO scheduler stopped")
			return
		case <-ticker.C:
			s.runEpochCycle(ctx)
		}
	}
}

func (s *Scheduler) runEpochCycle(ctx context.Context) {
	// Start epoch if needed
	if err := s.epochService.StartEpoch(ctx); err != nil {
		s.logger.Logf("ERROR failed to start epoch: %v", err)
	} else {
		s.logger.Logf("INFO successfully started epoch")
	}

	// Use vault address from configuration for subsidy distribution
	vaultId := s.config.Contracts.CollectionsVault
	if err := s.subsidyService.DistributeSubsidies(ctx, vaultId); err != nil {
		s.logger.Logf("ERROR failed to distribute subsidies: %v", err)
	} else {
		s.logger.Logf("INFO successfully distributed subsidies")
	}
}
