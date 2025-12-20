package openapi

import (
	"time"

	"github.com/nxlak/go-pvz/internal/domain/codes"
	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

func validateReturn(order *order_v1.Order) error {
	if order.Status != order_v1.OrderStatusISSUED && order.ExpiresAt.Value.Before(time.Now()) {
		return nil
	}
	return codes.ErrValidationFailed
}

func validateAccept(orderId, userId string, expiresAt time.Time) error {
	if orderId == "" {
		return codes.ErrValidationFailed
	}
	if userId == "" {
		return codes.ErrValidationFailed
	}
	if expiresAt.Before(time.Now()) {
		return codes.ErrValidationFailed
	}

	return nil
}
