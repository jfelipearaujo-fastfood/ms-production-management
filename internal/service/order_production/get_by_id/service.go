package get_by_id

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

func (s *Service) Handle(ctx context.Context, request GetOrderProductionByIdInput) (order_entity.Order, error) {
	if err := request.Validate(); err != nil {
		return order_entity.Order{}, err
	}

	order, err := s.repository.GetByID(ctx, request.OrderId)
	if err != nil {
		return order_entity.Order{}, err
	}

	order.RefreshStateTitle()

	return order, nil
}
