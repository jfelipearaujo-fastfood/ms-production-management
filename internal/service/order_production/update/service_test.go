package update

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	provider_mocks "github.com/jfelipearaujo-org/ms-production-management/internal/provider/mocks"
	repository_mocks "github.com/jfelipearaujo-org/ms-production-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should update order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		repository.On("Update", ctx, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Received",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		req := UpdateOrderProductionInput{
			OrderId: "order-id",
			State:   "Received",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to get the order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Received",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to update the order state", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{
				State: order_entity.Received,
			}, nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Completed",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should update order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		repository.On("Update", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Received",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
