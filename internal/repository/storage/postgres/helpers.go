package postgres

import (
	"github.com/jackc/pgx/v5"
	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

func scanOrders(rows pgx.Rows) ([]*order_v1.Order, error) {
	var orders []*order_v1.Order

	defer rows.Close()

	for rows.Next() {
		var o order_v1.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
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
