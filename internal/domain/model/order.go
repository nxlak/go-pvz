package model

import (
	"time"
)

type OrderStatus string

const (
	StatusAccepted OrderStatus = "ACCEPTED"

	StatusIssued OrderStatus = "ISSUED"

	StatusReturned OrderStatus = "RETURNED"

	StatusExpired OrderStatus = "EXPIRED"
)

type Order struct {
	Id         string      `json:"id"`
	UserId     string      `json:"user_id"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	ExpiresAt  time.Time   `json:"expires_at"`
	IssuedAt   time.Time   `json:"issued_at"`
	ReturnedAt time.Time   `json:"returned_at"`
}
