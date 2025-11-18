package service

import (
	"time"
)

type Service interface {
	AcceptOrder(orderId, userId string, expiresAt time.Time) error

	ReturnOrder(orderId string) error

	IssueOrder(userId string, orderIds []string) (map[string]error, error)
}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}
