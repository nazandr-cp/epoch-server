package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/andrey/epoch-server/internal/infra/config"
	"github.com/andrey/epoch-server/internal/services/epoch"
	"github.com/andrey/epoch-server/internal/services/subsidy"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduler_NewScheduler(t *testing.T) {
	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) (*epoch.StartEpochResponse, error) {
			return &epoch.StartEpochResponse{Status: "started"}, nil
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) (*epoch.ForceEndEpochResponse, error) {
			return &epoch.ForceEndEpochResponse{Status: "force_ended"}, nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &subsidy.ServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
			return &subsidy.SubsidyDistributionResponse{Status: "completed"}, nil
		},
	}

	logger := lgr.NoOp
	cfg := &config.Config{}
	cfg.Contracts.CollectionsVault = "0x1234567890123456789012345678901234567890"

	interval := 10 * time.Second

	scheduler := NewScheduler(mockEpochService, mockSubsidyService, interval, logger, cfg)

	require.NotNil(t, scheduler, "NewScheduler returned nil")
	require.NotNil(t, scheduler.epochService, "Scheduler epochService is nil")
	require.NotNil(t, scheduler.subsidyService, "Scheduler subsidyService is nil")
	require.NotNil(t, scheduler.logger, "Scheduler logger is nil")
	require.NotNil(t, scheduler.config, "Scheduler config is nil")

	assert.Equal(t, interval, scheduler.interval, "Interval mismatch")
}

func TestScheduler_runEpochCycle(t *testing.T) {
	epochStartCalled := false
	subsidyDistributeCalled := false

	mockEpochService := &epoch.ServiceMock{
		StartEpochFunc: func(ctx context.Context) (*epoch.StartEpochResponse, error) {
			epochStartCalled = true
			return &epoch.StartEpochResponse{Status: "started"}, nil
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) (*epoch.ForceEndEpochResponse, error) {
			return &epoch.ForceEndEpochResponse{Status: "force_ended"}, nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &subsidy.ServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
			subsidyDistributeCalled = true
			if vaultId != "0x1234567890123456789012345678901234567890" {
				t.Errorf("Expected vaultId 0x1234567890123456789012345678901234567890, got %s", vaultId)
			}
			return &subsidy.SubsidyDistributionResponse{Status: "completed"}, nil
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
		StartEpochFunc: func(ctx context.Context) (*epoch.StartEpochResponse, error) {
			epochStartCalled = true
			return nil, fmt.Errorf("epoch start error")
		},
		ForceEndEpochFunc: func(ctx context.Context, epochId uint64, vaultId string) (*epoch.ForceEndEpochResponse, error) {
			return &epoch.ForceEndEpochResponse{Status: "force_ended"}, nil
		},
		GetUserTotalEarnedFunc: func(ctx context.Context, userAddress, vaultId string) (*epoch.UserEarningsResponse, error) {
			return nil, nil
		},
	}

	mockSubsidyService := &subsidy.ServiceMock{
		DistributeSubsidiesFunc: func(ctx context.Context, vaultId string) (*subsidy.SubsidyDistributionResponse, error) {
			subsidyDistributeCalled = true
			return nil, fmt.Errorf("subsidy distribute error")
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
