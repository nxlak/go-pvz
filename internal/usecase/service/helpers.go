package service

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nxlak/go-pvz/internal/domain/model"
)

func appendOrder(order model.Order) error {
	f, err := os.OpenFile(OrdersFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(order)
	if err != nil {
		return err
	}

	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func deleteOrder(idx int) error {
	orders, err := readAllOrders()
	if err != nil {
		return err
	}

	orders = append(orders[:idx], orders[idx+1:]...)

	if err := writeAllOrders(orders); err != nil {
		return err
	}

	return nil
}

func validateReturn(order model.Order) error {
	if order.Status != model.StatusIssued && order.ExpiresAt.Before(time.Now()) {
		return nil
	}
	return fmt.Errorf("Can't delete this order")
}

func validateAccept(orderId, userId string, expiresAt time.Time) (model.Order, error) {
	orders, err := readAllOrders()
	if err != nil {
		return model.Order{}, err
	}

	for _, order := range orders {
		if order.Id == orderId {
			return model.Order{}, fmt.Errorf("Order already exists")
		}
	}

	if expiresAt.Before(time.Now()) {
		return model.Order{}, fmt.Errorf("Incorrect expire time")
	}

	order := model.Order{Id: orderId,
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Status:    model.StatusAccepted,
	}

	return order, nil
}
