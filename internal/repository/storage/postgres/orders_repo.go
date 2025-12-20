package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nxlak/go-pvz/internal/domain/codes"
	order "github.com/nxlak/go-pvz/internal/repository/storage"
	"github.com/nxlak/go-pvz/pkg/client/postgres"
	"github.com/nxlak/go-pvz/pkg/errs"
	"google.golang.org/protobuf/types/known/timestamppb"

	//order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
	orderV1 "github.com/nxlak/go-pvz/pkg/proto/order/v1"
)

type repository struct {
	client postgres.Client
}

func NewRepositoty(client postgres.Client) order.Repository {
	return &repository{client: client}
}

// gRPC VERSION

func (r *repository) Create(ctx context.Context, order *orderV1.Order) error {
	_, err := r.client.Exec(ctx, insertOrderSQL,
		order.Uuid,
		order.UserId,
		order.Status,
		order.CreatedAt.AsTime(),
		order.ExpiresAt.AsTime(),
		order.IssuedAt.AsTime(),
		order.ReturnedAt.AsTime(),
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
		"order_id", order.Uuid,
	)
}

func (r *repository) Update(ctx context.Context, order *orderV1.Order) error {
	tag, err := r.client.Exec(ctx, updateOrderSQL,
		order.Uuid,
		order.UserId,
		order.Status,
		order.ExpiresAt.AsTime(),
		order.IssuedAt.AsTime(),
		order.ReturnedAt.AsTime(),
	)
	if err != nil {
		return errs.Wrap(errs.CodeDatabaseError,
			"failed to update order", err, "order_id", order.Uuid)
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

func (r *repository) FindOne(ctx context.Context, id string) (*orderV1.Order, error) {
	var (
		o          orderV1.Order
		status     int32
		createdAt  time.Time
		expiresAt  time.Time
		issuedAt   pgtype.Timestamptz
		returnedAt pgtype.Timestamptz
	)

	if err := r.client.QueryRow(ctx, selectOrderSQL, id).Scan(
		&o.Uuid,
		&o.UserId,
		&status,
		&createdAt,
		&expiresAt,
		&issuedAt,
		&returnedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, codes.ErrOrderNotFound
		}
		return nil, errs.Wrap(errs.CodeDatabaseError,
			"failed to find order", err, "order_id", id)
	}

	o.Status = orderV1.OrderStatus(status)
	o.CreatedAt = timestamppb.New(createdAt)
	o.ExpiresAt = timestamppb.New(expiresAt)

	if issuedAt.Valid {
		o.IssuedAt = timestamppb.New(issuedAt.Time)
	}
	if returnedAt.Valid {
		o.ReturnedAt = timestamppb.New(returnedAt.Time)
	}

	return &o, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*orderV1.Order, error) {
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

func (r *repository) ListByUser(ctx context.Context, userId string) ([]*orderV1.Order, error) {
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

// HTTP VERSION
// func (r *repository) Create(ctx context.Context, order *order_v1.Order) error {
// 	_, err := r.client.Exec(ctx, insertOrderSQL,
// 		order.ID,
// 		order.UserID,
// 		order.Status,
// 		order.CreatedAt,
// 		order.ExpiresAt.Value,
// 		order.IssuedAt.Value,
// 		order.ReturnedAt.Value,
// 	)
// 	if err == nil {
// 		return nil
// 	}

// 	var pgErr *pgconn.PgError
// 	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
// 		return codes.ErrOrderAlreadyExists
// 	}

// 	return errs.Wrap(
// 		errs.CodeDatabaseError,
// 		"failed to insert order",
// 		err,
// 		"order_id", order.ID,
// 	)
// }

// func (r *repository) Update(ctx context.Context, order *orderV1.Order) error {
// 	tag, err := r.client.Exec(ctx, updateOrderSQL,
// 		order.ID,
// 		order.UserID,
// 		order.Status,
// 		order.ExpiresAt.Value,
// 		order.IssuedAt.Value,
// 		order.ReturnedAt.Value,
// 	)
// 	if err != nil {
// 		return errs.Wrap(errs.CodeDatabaseError,
// 			"failed to update order", err, "order_id", order.ID)
// 	}
// 	if tag.RowsAffected() == 0 {
// 		return codes.ErrOrderNotFound
// 	}
// 	return nil
// }

// func (r *repository) Delete(ctx context.Context, id string) error {
// 	tag, err := r.client.Exec(ctx, deleteOrderSQL, id)
// 	if err != nil {
// 		return errs.Wrap(errs.CodeDatabaseError,
// 			"failed to delete order", err, "order_id", id)
// 	}
// 	if tag.RowsAffected() == 0 {
// 		return codes.ErrOrderNotFound
// 	}
// 	return nil
// }

// func (r *repository) FindOne(ctx context.Context, id string) (*order_v1.Order, error) {
// 	var o order_v1.Order
// 	if err := r.client.QueryRow(ctx, selectOrderSQL, id).Scan(
// 		&o.ID,
// 		&o.UserID,
// 		&o.Status,
// 		&o.CreatedAt,
// 		&o.ExpiresAt.Value,
// 		&o.IssuedAt.Value,
// 		&o.ReturnedAt.Value,
// 	); err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, codes.ErrOrderNotFound
// 		}
// 		return nil, errs.Wrap(errs.CodeDatabaseError,
// 			"failed to find order", err, "order_id", id)
// 	}

// 	return &o, nil
// }

// func (r *repository) FindAll(ctx context.Context) ([]*order_v1.Order, error) {
// 	rows, err := r.client.Query(ctx, listAllSQL)
// 	if err != nil {
// 		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to list all orders", err)
// 	}
// 	defer rows.Close()

// 	orders, err := scanOrders(rows)
// 	if err != nil {
// 		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to scan orders", err)
// 	}

// 	return orders, nil
// }

// func (r *repository) ListByUser(ctx context.Context, userId string) ([]*order_v1.Order, error) {
// 	rows, err := r.client.Query(ctx, listByUserSQL, userId)
// 	if err != nil {
// 		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to list orders by user", err)
// 	}
// 	defer rows.Close()

// 	orders, err := scanOrders(rows)
// 	if err != nil {
// 		return nil, errs.Wrap(errs.CodeDatabaseError, "failed to scan orders", err)
// 	}

// 	return orders, nil
// }
