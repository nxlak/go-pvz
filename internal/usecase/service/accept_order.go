package service

import (
	"context"
	"time"

	"github.com/nxlak/go-pvz/internal/domain/model"
)

func (s *ServiceImpl) AcceptOrder(orderId, userId string, expiresAt time.Time) error {
	if err := validateAccept(orderId, userId, expiresAt); err != nil {
		return err
	}

	order := &model.Order{Id: orderId,
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Status:    model.StatusAccepted,
	}

	if err := s.orderRepo.Create(context.TODO(), order); err != nil {
		return err
	}

	return nil
}
