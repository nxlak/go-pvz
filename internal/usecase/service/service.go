package service

import (
	"time"

	order "github.com/nxlak/go-pvz/internal/repository/storage"
)

type Service interface {
	AcceptOrder(orderId, userId string, expiresAt time.Time) error

	ReturnOrder(orderId string) error

	IssueOrder(userId string, orderIds []string) (map[string]error, error)

	ReturnOrdersByUser(userId string, orderIds []string) (map[string]error, error)
}

type ServiceImpl struct {
	orderRepo order.Repository
}

func NewService(orderRepo order.Repository) *ServiceImpl {
	return &ServiceImpl{orderRepo: orderRepo}
}
