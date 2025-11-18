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
	Id          string      `json:"id"`
	UserId      string      `json:"user_id"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
	CompletedAt time.Time   `json:"completed_at"`
	Status      OrderStatus `json:"status"`
}
