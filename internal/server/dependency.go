package server

import (
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_id"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	OrderProductionRepository repository.OrderProductionRepository

	GetOrderProductionById service.GetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput]
}
