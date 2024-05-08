package order_production

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
)

type OrderProductionRepository struct {
	conn *sql.DB
}

func NewOrderProductionRepository(conn *sql.DB) *OrderProductionRepository {
	return &OrderProductionRepository{
		conn: conn,
	}
}

func (r *OrderProductionRepository) Create(ctx context.Context, order *order_entity.Order) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sql, params, err := goqu.
		Insert("orders").
		Cols("order_id", "state", "state_updated_at", "created_at", "updated_at").
		Vals(
			goqu.Vals{
				order.Id,
				order.State,
				order.StateUpdatedAt,
				order.CreatedAt,
				order.UpdatedAt,
			},
		).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sql, params...)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return errTx
		}
		return err
	}

	for _, item := range order.Items {
		sql, params, err := goqu.
			Insert("order_items").
			Cols("id", "order_id", "name", "quantity").
			Vals(
				goqu.Vals{
					item.Id,
					order.Id,
					item.Name,
					item.Quantity,
				},
			).
			ToSQL()
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, sql, params...)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return errTx
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderProductionRepository) GetByID(ctx context.Context, id string) (order_entity.Order, error) {
	var order order_entity.Order

	order.Items = make([]order_entity.Item, 0)

	sql, params, err := goqu.
		From("orders").
		Select("order_id", "state", "state_updated_at", "created_at", "updated_at").
		Where(goqu.C("order_id").Eq(id)).
		ToSQL()
	if err != nil {
		return order_entity.Order{}, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return order_entity.Order{}, err
	}
	defer statement.Close()

	for statement.Next() {
		if err := statement.Scan(
			&order.Id,
			&order.State,
			&order.StateUpdatedAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return order_entity.Order{}, err
		}
	}

	if order.Id == "" {
		return order_entity.Order{}, custom_error.ErrOrderNotFound
	}

	sql, params, err = goqu.
		From("order_items").
		Select("id", "name", "quantity").
		Where(goqu.C("order_id").Eq(order.Id)).
		ToSQL()
	if err != nil {
		return order_entity.Order{}, err
	}

	rows, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return order_entity.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item order_entity.Item

		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Quantity,
		); err != nil {
			return order_entity.Order{}, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderProductionRepository) GetByState(ctx context.Context, state order_entity.OrderState) ([]order_entity.Order, error) {
	var orders []order_entity.Order

	sql, params, err := goqu.
		From("orders").
		Select("order_id", "state", "state_updated_at", "created_at", "updated_at").
		Where(goqu.C("state").Eq(state)).
		Order(goqu.C("created_at").Desc()).
		ToSQL()
	if err != nil {
		return orders, err
	}

	statement, err := r.conn.QueryContext(ctx, sql, params...)
	if err != nil {
		return orders, err
	}
	defer statement.Close()

	for statement.Next() {
		var order order_entity.Order

		if err := statement.Scan(
			&order.Id,
			&order.State,
			&order.StateUpdatedAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return orders, err
		}

		order.RefreshStateTitle()

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderProductionRepository) Update(ctx context.Context, order *order_entity.Order) error {
	sql, params, err := goqu.
		Update("orders").
		Set(goqu.Record{
			"state":            order.State,
			"state_updated_at": order.StateUpdatedAt,
			"updated_at":       order.UpdatedAt,
		}).
		Where(goqu.C("order_id").Eq(order.Id)).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.conn.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}

	return nil
}
