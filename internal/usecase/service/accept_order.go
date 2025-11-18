package service

import (
	"time"
)

func (s *ServiceImpl) AcceptOrder(orderId, userId string, expiresAt time.Time) error {
	order, err := validateAccept(orderId, userId, expiresAt)
	if err != nil {
		return err
	}

	if err := appendOrder(order); err != nil {
		return err
	}

	return nil
}
