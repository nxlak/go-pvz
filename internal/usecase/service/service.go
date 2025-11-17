package service

import (
	"time"
)

type Service interface {
	AcceptOrder(orderId, UserId string, expiresAt time.Time) error

	// ...

}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}
