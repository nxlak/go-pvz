package server

import (
	"time"

	"github.com/nxlak/go-pvz/internal/domain/codes"
	orderV1 "github.com/nxlak/go-pvz/pkg/proto/order/v1"
)

func validateReturn(order *orderV1.Order) error {
	if order.Status != orderV1.OrderStatus_ORDER_STATUS_ISSUED && order.ExpiresAt.AsTime().Before(time.Now()) {
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
