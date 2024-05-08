package service

import (
	"context"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
)

type CreateOrderProductionService[T any] interface {
	Handle(ctx context.Context, request T) (*order_entity.Order, error)
}

type GetOrderProductionByIdService[T any] interface {
	Handle(ctx context.Context, request T) (order_entity.Order, error)
}

type GetOrderProductionByStateService[T any] interface {
	Handle(ctx context.Context, request T) ([]order_entity.Order, error)
}

type UpdateOrderProductionService[T any] interface {
	Handle(ctx context.Context, request T) (*order_entity.Order, error)
}
