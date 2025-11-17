package service

import (
	"errors"
	"time"
)

func (s *ServiceImpl) AcceptOrder(orderId, userId string, expiresAt time.Time) error {
	if err := validateAccept(orderId, userId, expiresAt); err != nil {
		return errors.New("Invalid data")
	}

	return nil
}
