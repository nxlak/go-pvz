package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nxlak/go-pvz/internal/domain/model"
	order "github.com/nxlak/go-pvz/internal/repository/storage"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
)

type repository struct {
	client postgres.Client
}

const (
	insertOrderSQL = `
		INSERT INTO orders (
			id, user_id, status, created_at, 
			expires_at, issued_at, returned_at
		) VALUES (
			$1,$2,$3,$4,
			$5,$6,$7
		);`

	updateOrderSQL = `
		UPDATE orders SET
			user_id      = $2,
			status       = $3,
			expires_at   = $4,
			issued_at    = $5,
			returned_at  = $6,
			package      = $7,
			weight       = $8,
			price        = $9,
			total_price  = $10
		WHERE id = $1;`

	selectOrderSQL = `
		SELECT
			id, user_id, status, expires_at,
			issued_at, returned_at, created_at,
			package, weight, price, total_price
		FROM orders
		WHERE id = $1;`

	deleteOrderSQL = `DELETE FROM orders WHERE id = $1;`

	listAllSQL = `
		SELECT
			id, user_id, status, expires_at,
			issued_at, returned_at, created_at,
			package, weight, price, total_price
		FROM orders
		ORDER BY created_at ASC, id ASC;`
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	_, err := r.client.Exec(ctx, insertOrderSQL,
		order.Id,
		order.UserId,
		order.Status,
		order.CreatedAt,
		order.ExpiresAt,
		order.IssuedAt,
		order.ReturnedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newError := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			fmt.Println(newError)
			return nil
		}
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Exec(ctx, deleteOrderSQL, id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newError := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			fmt.Println(newError)
			return nil
		}
		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) (o []model.Order, err error) {
	panic("unimplemented")
}

func (r *repository) FindOne(ctx context.Context, id string) (model.Order, error) {
	panic("unimplemented")
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	panic("unimplemented")
}

func NewRepositoty(client postgres.Client) order.Repository {
	return &repository{client: client}
}
