package update

import (
	"context"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
)

type Service struct {
	repository   repository.OrderProductionRepository
	timeProvider provider.TimeProvider
}

func NewService(
	repository repository.OrderProductionRepository,
	timeProvider provider.TimeProvider,
) *Service {
	return &Service{
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (s *Service) Handle(ctx context.Context, request UpdateOrderProductionInput) (*order_entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	order, err := s.repository.GetByID(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}

	newState := order_entity.NewOrderState(request.State)

	if err := order.UpdateState(newState, s.timeProvider.GetTime()); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
