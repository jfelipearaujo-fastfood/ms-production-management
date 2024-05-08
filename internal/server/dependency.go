package server

import (
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_id"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_state"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	OrderProductionRepository repository.OrderProductionRepository

	GetOrderProductionById    service.GetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput]
	GetOrderProductionByState service.GetOrderProductionByStateService[get_by_state.GetOrderProductionByStateInput]
}
