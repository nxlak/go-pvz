package order

import (
	"context"

	"github.com/nxlak/go-pvz/internal/domain/model"
)

type Repository interface {
	Create(ctx context.Context, order *model.Order) error
	FindAll(ctx context.Context) (o []*model.Order, err error)
	FindOne(ctx context.Context, id string) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userId string) ([]*model.Order, error)
}
