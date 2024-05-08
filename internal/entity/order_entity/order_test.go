package order_entity

import (
	"testing"
	"time"

	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	t.Run("Should create a new order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		// Act
		res := NewOrder("customer_id", now)

		// Assert
		assert.NotEmpty(t, res.Id)
		assert.Equal(t, Received, res.State)
		assert.Equal(t, now, res.StateUpdatedAt)
		assert.Empty(t, res.Items)
		assert.Equal(t, now, res.CreatedAt)
		assert.Equal(t, now, res.UpdatedAt)
	})

	t.Run("Should add an item to the order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		expectedItem := Item{
			Id:       "item_id",
			Name:     "name",
			Quantity: 1,
		}

		order := NewOrder("customer_id", now)

		// Act
		err := order.AddItem(NewItem("item_id", "name", 1), now)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, order.Items, 1)
		assert.Contains(t, order.Items, expectedItem)
	})

	t.Run("Should update the state of the order", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Processing, now)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, Processing, order.State)
		assert.Equal(t, now, order.StateUpdatedAt)
		assert.Equal(t, now, order.UpdatedAt)
	})

	t.Run("Should return an error when trying to update the state to an invalid state", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Completed, now)

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should not update the state if is the same state", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Received, now)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, Received, order.State)
		assert.Equal(t, past, order.StateUpdatedAt)
		assert.Equal(t, past, order.UpdatedAt)
	})

	t.Run("Should refresh the state title", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		order.RefreshStateTitle()

		// Assert
		assert.Equal(t, "Received", order.StateTitle)
	})

	t.Run("Should return true if the order is already completed", func(t *testing.T) {
		// Arrange
		states := []OrderState{Delivered, Cancelled}

		for _, state := range states {
			now := time.Now()

			order := NewOrder("customer_id", now)
			order.State = state

			// Act
			res := order.IsCompleted()

			// Assert
			assert.True(t, res)
		}
	})

	t.Run("Should return false if the order is not completed", func(t *testing.T) {
		// Arrange
		states := []OrderState{
			None,
			Received,
			Processing,
			Completed,
		}

		for _, state := range states {
			now := time.Now()

			order := NewOrder("customer_id", now)
			order.State = state

			// Act
			res := order.IsCompleted()

			// Assert
			assert.False(t, res)
		}
	})

	t.Run("Should return an error when trying to add an item that already exists", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		order.Items = append(order.Items, NewItem("item_id", "name", 1))

		// Act
		err := order.AddItem(NewItem("item_id", "name", 1), now)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrOrderItemAlreadyExists)
	})

	t.Run("Should return true if the order has items", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		order.Items = append(order.Items, NewItem("item_id", "name", 1))

		// Act
		res := order.HasItems()

		// Assert
		assert.True(t, res)
	})
}
