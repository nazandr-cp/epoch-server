package service

import (
	"context"
	"fmt"

	"github.com/andrey/epoch-server/internal/clients/contract"
	"github.com/andrey/epoch-server/internal/clients/epoch"
	"github.com/andrey/epoch-server/internal/clients/graph"
	"github.com/andrey/epoch-server/internal/clients/storage"
	"github.com/andrey/epoch-server/internal/clients/subsidizer"
	"github.com/go-pkgz/lgr"
)

type GraphClient interface {
	QueryUsers(ctx context.Context) ([]graph.User, error)
	QueryEligibility(ctx context.Context, epochID string) ([]graph.Eligibility, error)
	ExecuteQuery(ctx context.Context, request graph.GraphQLRequest, response interface{}) error
	ExecutePaginatedQuery(ctx context.Context, queryTemplate string, variables map[string]interface{}, entityField string, response interface{}) error
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

func (s *Service) DistributeSubsidies(ctx context.Context, vaultId string) error {
	epochClient := epoch.NewClient(s.logger)
	subsidizerClient := subsidizer.NewClient(s.logger)
	storageClient := storage.NewClient(s.logger)

	lazyDistributor := NewLazyDistributor(
		s.graphClient,
		epochClient,
		subsidizerClient,
		storageClient,
		s.logger,
	)

	if err := lazyDistributor.Run(ctx, vaultId); err != nil {
		return fmt.Errorf("failed to run lazy distributor for vault %s: %w", vaultId, err)
	}

	return nil
}
