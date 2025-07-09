package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/go-pkgz/lgr"
)

func TestScheduler_NewScheduler(t *testing.T) {
	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) error {
			return nil
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) error {
			return nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &SubsidyServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) error {
			return nil
		},
	}

	logger := lgr.NoOp
	cfg := &config.Config{}
	cfg.Contracts.CollectionsVault = "0x1234567890123456789012345678901234567890"

	interval := 10 * time.Second

	scheduler := NewScheduler(mockEpochService, mockSubsidyService, interval, logger, cfg)

	if scheduler == nil {
		t.Error("NewScheduler returned nil")
	}

	if scheduler.epochService == nil {
		t.Error("Scheduler epochService is nil")
	}

	if scheduler.subsidyService == nil {
		t.Error("Scheduler subsidyService is nil")
	}

	if scheduler.logger == nil {
		t.Error("Scheduler logger is nil")
	}

	if scheduler.interval != interval {
		t.Errorf("Expected interval %v, got %v", interval, scheduler.interval)
	}

	if scheduler.config == nil {
		t.Error("Scheduler config is nil")
	}
}

func TestScheduler_runEpochCycle(t *testing.T) {
	epochStartCalled := false
	subsidyDistributeCalled := false

	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) error {
			epochStartCalled = true
			return nil
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) error {
			return nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &SubsidyServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) error {
			subsidyDistributeCalled = true
			if vaultId != "0x1234567890123456789012345678901234567890" {
				t.Errorf("Expected vaultId 0x1234567890123456789012345678901234567890, got %s", vaultId)
			}
			return nil
		},
	}

	logger := lgr.NoOp
	cfg := &config.Config{}
	cfg.Contracts.CollectionsVault = "0x1234567890123456789012345678901234567890"

	interval := 10 * time.Second

	scheduler := NewScheduler(mockEpochService, mockSubsidyService, interval, logger, cfg)

	ctx := context.Background()
	scheduler.runEpochCycle(ctx)

	if !epochStartCalled {
		t.Error("Expected StartEpoch to be called")
	}

	if !subsidyDistributeCalled {
		t.Error("Expected DistributeSubsidies to be called")
	}
}

func TestScheduler_runEpochCycle_WithErrors(t *testing.T) {
	epochStartCalled := false
	subsidyDistributeCalled := false

	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) error {
			epochStartCalled = true
			return fmt.Errorf("epoch start error")
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) error {
			return nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &SubsidyServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) error {
			subsidyDistributeCalled = true
			return fmt.Errorf("subsidy distribute error")
		},
	}

	logger := lgr.NoOp
	cfg := &config.Config{}
	cfg.Contracts.CollectionsVault = "0x1234567890123456789012345678901234567890"

	interval := 10 * time.Second

	scheduler := NewScheduler(mockEpochService, mockSubsidyService, interval, logger, cfg)

	ctx := context.Background()
	scheduler.runEpochCycle(ctx)

	if !epochStartCalled {
		t.Error("Expected StartEpoch to be called")
	}

	if !subsidyDistributeCalled {
		t.Error("Expected DistributeSubsidies to be called")
	}
}
