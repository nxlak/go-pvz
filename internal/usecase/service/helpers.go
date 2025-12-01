package service

import (
	"time"

	"github.com/nxlak/go-pvz/internal/domain/codes"
	"github.com/nxlak/go-pvz/internal/domain/model"
)

func validateReturn(order *model.Order) error {
	if order.Status != model.StatusIssued && order.ExpiresAt.Before(time.Now()) {
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
