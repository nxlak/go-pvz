package model

import (
	"time"
)

type OrderStatus string

const (
	StatusAccepted OrderStatus = "ACCEPTED"

	StatusCompleted OrderStatus = "COMPLETED" // заказ выдан человеку

	StatusReturned OrderStatus = "RETURNED"

	StatusExpired OrderStatus = "EXPIRED"
)

type Order struct {
	Id          string
	UserId      string
	ExpiresAt   time.Time
	CompletedAt time.Time
	Status      OrderStatus
}
