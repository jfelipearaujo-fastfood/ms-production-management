package get_by_state

import (
	"context"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
)

type Service struct {
	repository repository.OrderProductionRepository
}

func NewService(
	repository repository.OrderProductionRepository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Handle(ctx context.Context, request GetOrderProductionByStateInput) ([]order_entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	state := order_entity.NewOrderState(request.State)

	orders, err := s.repository.GetByState(ctx, state)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
