package order_production

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Should create a new order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?orders(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO (.+)?order_items(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		repo := NewOrderProductionRepository(db)

		order := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		order.AddItem(orderItem, time.Now())

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when try to begin the transaction", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin().
			WillReturnError(assert.AnError)

		repo := NewOrderProductionRepository(db)

		order := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		order.AddItem(orderItem, time.Now())

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return error when order insert fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?orders(.+)?").
			WillReturnError(assert.AnError)

		mock.ExpectRollback()

		repo := NewOrderProductionRepository(db)

		order := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		order.AddItem(orderItem, time.Now())

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return error when order items insert fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO (.+)?orders(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO (.+)?order_items(.+)?").
			WillReturnError(assert.AnError)

		mock.ExpectRollback()

		repo := NewOrderProductionRepository(db)

		order := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		order.AddItem(orderItem, time.Now())

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
	})
}
