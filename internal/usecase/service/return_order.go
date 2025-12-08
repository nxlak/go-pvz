package service

// import "context"

// func (s *ServiceImpl) ReturnOrder(orderId string) error {
// 	order, err := s.orderRepo.FindOne(context.TODO(), orderId)
// 	if err != nil {
// 		return err
// 	}

// 	if err := validateReturn(order); err != nil {
// 		return err
// 	}

// 	if err := s.orderRepo.Delete(context.TODO(), orderId); err != nil {
// 		return err
// 	}

// 	return nil
// }
