package create

import (
	"context"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
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

func (s *Service) Handle(ctx context.Context, request CreateOrderProductionInput) (*order_entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	exists, err := s.repository.GetByID(ctx, request.OrderId)
	if err != nil && err != custom_error.ErrOrderNotFound {
		return nil, err
	}

	if exists.Exists() {
		return nil, custom_error.ErrOrderAlreadyExists
	}

	order := order_entity.NewOrder(request.OrderId, s.timeProvider.GetTime())

	for _, item := range request.Items {
		orderItem := order_entity.NewItem(item.Id, item.Name, item.Quantity)

		if err := order.AddItem(orderItem, s.timeProvider.GetTime()); err != nil {
			return nil, err
		}
	}

	if err := s.repository.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
