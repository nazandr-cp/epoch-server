package service

import (
	"context"
	"fmt"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/go-pkgz/lgr"
)

type GraphClient interface {
	QueryUsers(ctx context.Context) ([]graph.User, error)
	QueryEligibility(ctx context.Context, epochID string) ([]graph.Eligibility, error)
}

type ContractClient interface {
	StartEpoch(ctx context.Context, epochID string) error
	DistributeSubsidies(ctx context.Context, epochID string) error
}

type Service struct {
	graphClient    GraphClient
	contractClient ContractClient
	logger         lgr.L
}

func NewService(graphClient *graph.Client, contractClient *contract.Client, logger lgr.L) *Service {
	return &Service{
		graphClient:    graphClient,
		contractClient: contractClient,
		logger:         logger,
	}
}

func (s *Service) StartEpoch(ctx context.Context, epochID string) error {
	users, err := s.graphClient.QueryUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}

	eligibilities, err := s.graphClient.QueryEligibility(ctx, epochID)
	if err != nil {
		return fmt.Errorf("failed to query eligibility for epoch %s: %w", epochID, err)
	}

	s.logger.Logf("INFO found %d users and %d eligibilities for epoch %s", len(users), len(eligibilities), epochID)

	if err := s.contractClient.StartEpoch(ctx, epochID); err != nil {
		return fmt.Errorf("failed to start epoch %s: %w", epochID, err)
	}

	return nil
}

func (s *Service) DistributeSubsidies(ctx context.Context, epochID string) error {
	eligibilities, err := s.graphClient.QueryEligibility(ctx, epochID)
	if err != nil {
		return fmt.Errorf("failed to query eligibility for epoch %s: %w", epochID, err)
	}

	eligibleCount := 0
	for _, eligibility := range eligibilities {
		if eligibility.IsEligible {
			eligibleCount++
		}
	}

	s.logger.Logf("INFO found %d eligible users for epoch %s", eligibleCount, epochID)

	if err := s.contractClient.DistributeSubsidies(ctx, epochID); err != nil {
		return fmt.Errorf("failed to distribute subsidies for epoch %s: %w", epochID, err)
	}

	return nil
}
