package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nxlak/go-pvz/internal/domain/model"
)

const OrdersFilePath = "orders.txt"

func readAllOrders() ([]model.Order, error) {
	file, err := os.Open(OrdersFilePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var orders []model.Order
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var order model.Order
		if err := json.Unmarshal(line, &order); err != nil {
			return nil, fmt.Errorf("failed to parse order line: %w", err)
		}

		orders = append(orders, order)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

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

func validateAccept(orderId, userId string, expiresAt time.Time) (model.Order, error) {
	orders, err := readAllOrders()
	if err != nil {
		return model.Order{}, err
	}

	for _, order := range orders {
		if order.Id == orderId {
			return model.Order{}, fmt.Errorf("Order already exists!")
		}
	}

	now := time.Now()
	if expiresAt.Before(now) {
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
