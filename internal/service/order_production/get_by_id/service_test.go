package get_by_id

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		service := NewService(repository)

		req := GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, order)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderProductionRepository(t)

		service := NewService(repository)

		req := GetOrderProductionByIdInput{
			OrderId: "123",
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
		repository.AssertExpectations(t)
	})
}
