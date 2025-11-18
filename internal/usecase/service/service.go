package service

import (
	"time"
)

type Service interface {
	AcceptOrder(orderId, UserId string, expiresAt time.Time) error

	ReturnOrder(orderId string) error
}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}
