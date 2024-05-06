package order_production

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
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
		Cols("id", "state", "state_updated_at", "created_at", "updated_at").
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
