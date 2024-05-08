package get_by_state

import (
	"context"
	"testing"

	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return orders", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		repository.On("GetByState", ctx, mock.Anything).
			Return([]order_entity.Order{}, nil).
			Once()

		service := NewService(repository)

		req := GetOrderProductionByStateInput{
			State: "Received",
		}

		// Act
		orders, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		repository.On("GetByState", ctx, mock.Anything).
			Return([]order_entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetOrderProductionByStateInput{
			State: "Received",
		}

		// Act
		orders, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, orders)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		service := NewService(repository)

		req := GetOrderProductionByStateInput{
			State: "123",
		}

		// Act
		orders, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, orders)
		repository.AssertExpectations(t)
	})
}
