package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nxlak/go-pvz/internal/domain/model"
)

const OrdersFilePath = "orders.txt"

func writeAllOrders(orders []model.Order) error {
	dir := filepath.Dir(OrdersFilePath)
	tmp, err := os.CreateTemp(dir, "orders-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()

	enc := json.NewEncoder(tmp)
	for _, o := range orders {
		if err := enc.Encode(o); err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
			return fmt.Errorf("encode order %s: %w", o.Id, err)
		}
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("sync temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpName, OrdersFilePath); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

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
