package postgres

import (
	"github.com/jackc/pgx/v5"
	"github.com/nxlak/go-pvz/internal/domain/model"
)

func scanOrders(rows pgx.Rows) ([]*model.Order, error) {
	var orders []*model.Order

	defer rows.Close()

	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.Id,
			&o.UserId,
			&o.Status,
			&o.CreatedAt,
			&o.ExpiresAt,
			&o.IssuedAt,
			&o.ReturnedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}
