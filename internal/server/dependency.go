package server

import (
	"github.com/jfelipearaujo-org/ms-production-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	OrderProductionRepository repository.OrderProductionRepository
}
