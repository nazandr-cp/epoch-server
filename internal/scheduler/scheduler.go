package scheduler

import (
	"context"
	"time"

	"github.com/andrey/epoch-server/internal/service"
	"github.com/go-pkgz/lgr"
)

type Scheduler struct {
	service  *service.Service
	logger   lgr.L
	interval time.Duration
}

func NewScheduler(interval time.Duration, svc *service.Service, logger lgr.L) *Scheduler {
	return &Scheduler{
		service:  svc,
		logger:   logger,
		interval: interval,
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
	epochID := ""

	if err := s.service.StartEpoch(ctx, epochID); err != nil {
		s.logger.Logf("ERROR failed to start epoch: %v", err)
	} else {
		s.logger.Logf("INFO successfully started epoch")
	}

	if err := s.service.DistributeSubsidies(ctx, epochID); err != nil {
		s.logger.Logf("ERROR failed to distribute subsidies: %v", err)
	} else {
		s.logger.Logf("INFO successfully distributed subsidies")
	}
}
