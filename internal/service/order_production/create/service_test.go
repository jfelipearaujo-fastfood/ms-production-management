package create

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
	t.Run("Should create order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		repository.On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		req := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
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

		req := CreateOrderProductionInput{
			OrderId: "order-id",
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		repository.On("Create", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		req := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to add item fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(3)

		service := NewService(repository, timeProvider)

		itemId := uuid.NewString()

		req := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       itemId,
					Name:     "Test",
					Quantity: 1,
				},
				{
					Id:       itemId,
					Name:     "Test",
					Quantity: 1,
				},
			},
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should not create order when order already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{
				Id: uuid.NewString(),
			}, nil).
			Once()

		service := NewService(repository, timeProvider)

		req := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
		}

		// Act
		order, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when could not get order by ID", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mocks.NewMockOrderProductionRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository, timeProvider)

		req := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
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
