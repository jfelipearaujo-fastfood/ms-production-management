package order_production

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
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
		err = order.AddItem(orderItem, time.Now())
		assert.NoError(t, err)

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
		err = order.AddItem(orderItem, time.Now())
		assert.NoError(t, err)

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
		err = order.AddItem(orderItem, time.Now())
		assert.NoError(t, err)

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
		err = order.AddItem(orderItem, time.Now())
		assert.NoError(t, err)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("Should return order without items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, expectedOrder.State, expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?order_items(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "name", "quantity"}))

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, order)
	})

	t.Run("Should return order with items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			now,
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		err = expectedOrder.AddItem(orderItem, now)
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, expectedOrder.State, expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?order_items(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "quantity"}).
				AddRow(orderItem.Id, orderItem.Name, orderItem.Quantity))

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, order)
	})

	t.Run("Should return scan error when find the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, "abc", expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, "id")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
	})

	t.Run("Should return scan error when find the order items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, expectedOrder.State, expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?order_items(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "quantity"}).
				AddRow("id", "name", "quantity"))

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
	})

	t.Run("Should return error when find the order items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			now,
		)
		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		err = expectedOrder.AddItem(orderItem, now)
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, expectedOrder.State, expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?order_items(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
	})

	t.Run("Should return error when find the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, order)
	})

	t.Run("Should return error when order is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}))

		repo := NewOrderProductionRepository(db)

		// Act
		order, err := repo.GetByID(ctx, expectedOrder.Id)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, custom_error.ErrOrderNotFound, err)
		assert.Empty(t, order)
	})
}

func TestGetByState(t *testing.T) {
	t.Run("Should return orders with state", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			now,
		)

		orderItem := order_entity.NewItem(
			uuid.NewString(),
			"Item",
			1,
		)
		err = expectedOrder.AddItem(orderItem, now)
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, expectedOrder.State, expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		mock.ExpectQuery("SELECT (.+)?order_items(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "quantity"}).
				AddRow(orderItem.Id, orderItem.Name, orderItem.Quantity))

		repo := NewOrderProductionRepository(db)

		// Act
		orders, err := repo.GetByState(ctx, expectedOrder.State)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, orders, 1)
	})

	t.Run("Should return empty if no orders were found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}))

		repo := NewOrderProductionRepository(db)

		// Act
		orders, err := repo.GetByState(ctx, expectedOrder.State)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})

	t.Run("Should return error when find the orders", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderProductionRepository(db)

		// Act
		orders, err := repo.GetByState(ctx, expectedOrder.State)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, orders)
	})

	t.Run("Should return error when scan fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectQuery("SELECT (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"id", "state", "state_updated_at", "created_at", "updated_at"}).
				AddRow(expectedOrder.Id, "abc", expectedOrder.StateUpdatedAt, expectedOrder.CreatedAt, expectedOrder.UpdatedAt))

		repo := NewOrderProductionRepository(db)

		// Act
		orders, err := repo.GetByState(ctx, expectedOrder.State)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, orders)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Should update the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			now,
		)

		mock.ExpectExec("UPDATE (.+)?orders(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewOrderProductionRepository(db)

		// Act
		err = repo.Update(ctx, &expectedOrder)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when update fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		expectedOrder := order_entity.NewOrder(
			uuid.NewString(),
			time.Now(),
		)

		mock.ExpectExec("UPDATE (.+)?orders(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderProductionRepository(db)

		// Act
		err = repo.Update(ctx, &expectedOrder)

		// Assert
		assert.Error(t, err)
	})
}
