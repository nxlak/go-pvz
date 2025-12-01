package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nxlak/go-pvz/internal/domain/codes"
	"github.com/nxlak/go-pvz/internal/domain/model"
	order "github.com/nxlak/go-pvz/internal/repository/storage"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
	"github.com/nxlak/go-pvz/pkg/errs"
)

type repository struct {
	client postgres.Client
}

func NewRepositoty(client postgres.Client) order.Repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, order *model.Order) error {
	_, err := r.client.Exec(ctx, insertOrderSQL,
		order.Id,
		order.UserId,
		order.Status,
		order.CreatedAt,
		order.ExpiresAt,
		order.IssuedAt,
		order.ReturnedAt,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return codes.ErrOrderAlreadyExists
	}

	return errs.Wrap(
		errs.CodeDatabaseError,
		"failed to insert order",
		err,
		"order_id", order.Id,
	)
}

func (r *repository) Update(ctx context.Context, order *model.Order) error {
	tag, err := r.client.Exec(ctx, updateOrderSQL,
		order.Id,
		order.UserId,
		order.Status,
		order.CreatedAt,
		order.ExpiresAt,
		order.IssuedAt,
		order.ReturnedAt,
	)
	if err != nil {
		return errs.Wrap(errs.CodeDatabaseError,
			"failed to update order", err, "order_id", order.Id)
	}
	if tag.RowsAffected() == 0 {
		return codes.ErrOrderNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	tag, err := r.client.Exec(ctx, deleteOrderSQL, id)
	if err != nil {
		return errs.Wrap(errs.CodeDatabaseError,
			"failed to delete order", err, "order_id", id)
	}
	if tag.RowsAffected() == 0 {
		return codes.ErrOrderNotFound
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id string) (*model.Order, error) {
	var o model.Order
	if err := r.client.QueryRow(ctx, selectOrderSQL, id).Scan(
		&o.Id,
		&o.UserId,
		&o.Status,
		&o.CreatedAt,
		&o.ExpiresAt,
		&o.IssuedAt,
		&o.ReturnedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, codes.ErrOrderNotFound
		}
		return nil, errs.Wrap(errs.CodeDatabaseError,
			"failed to find order", err, "order_id", id)
	}

	return &o, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.client.Query(ctx, listAllSQL)
	if err != nil {
		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to list all orders", err)
	}
	defer rows.Close()

	orders, err := scanOrders(rows)
	if err != nil {
		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to scan orders", err)
	}

	return orders, nil
}

func (r *repository) ListByUser(ctx context.Context, userId string) ([]*model.Order, error) {
	rows, err := r.client.Query(ctx, listByUserSQL, userId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to list orders by user", err)
	}
	defer rows.Close()

	orders, err := scanOrders(rows)
	if err != nil {
		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to scan orders", err)
	}

	return orders, nil
}
