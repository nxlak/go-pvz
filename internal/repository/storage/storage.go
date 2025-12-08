package order

import (
	"context"

	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

type Repository interface {
	Create(ctx context.Context, order *order_v1.Order) error
	FindAll(ctx context.Context) (o []*order_v1.Order, err error)
	FindOne(ctx context.Context, id string) (*order_v1.Order, error)
	Update(ctx context.Context, order *order_v1.Order) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userId string) ([]*order_v1.Order, error)
}
