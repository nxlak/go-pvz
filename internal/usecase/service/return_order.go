package service

import (
	"fmt"
)

func (s *ServiceImpl) ReturnOrder(orderId string) error {
	orders, err := readAllOrders()
	if err != nil {
		return err
	}

	idx := -1
	for i, o := range orders {
		if o.Id == orderId {
			idx = i
		}
	}

	if idx == -1 {
		return fmt.Errorf("order wasnt found")
	}
	order := orders[idx]

	if err := validateReturn(order); err != nil {
		return err
	}

	if err := deleteOrder(idx); err != nil {
		return err
	}

	return nil
}
