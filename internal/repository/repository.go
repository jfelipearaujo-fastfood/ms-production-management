package repository

import (
	"context"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
)

type OrderProductionRepository interface {
	Create(ctx context.Context, order *order_entity.Order) error
	GetByID(ctx context.Context, id string) (order_entity.Order, error)
	GetByState(ctx context.Context, state order_entity.OrderState) ([]order_entity.Order, error)
	Update(ctx context.Context, order *order_entity.Order) error
}
