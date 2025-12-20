package server

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	order "github.com/nxlak/go-pvz/internal/repository/storage"
	orderV1 "github.com/nxlak/go-pvz/pkg/proto/order/v1"
)

type OrderService struct {
	orderV1.UnimplementedOrderServiceServer
	storage order.Repository
}

func NewOrderService(storage order.Repository) *OrderService {
	return &OrderService{
		storage: storage,
	}
}

func (s *OrderService) GetOrderById(ctx context.Context, req *orderV1.GetOrderByIdRequest) (*orderV1.GetOrderByIdResponse, error) {
	order, err := s.storage.FindOne(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &orderV1.GetOrderByIdResponse{
		Order: order,
	}, nil
}

func (s *OrderService) AcceptOrder(ctx context.Context, req *orderV1.AcceptOrderRequest) (*orderV1.OrderResponse, error) {
	if err := validateAccept(req.Uuid, req.UserId, req.ExpiresAt.AsTime()); err != nil {
		return nil, err
	}

	order := orderV1.Order{
		Uuid:      req.Uuid,
		UserId:    req.UserId,
		CreatedAt: timestamppb.New(time.Now()),
		ExpiresAt: req.ExpiresAt,
		Status:    orderV1.OrderStatus_ORDER_STATUS_ACCEPTED,
	}

	if err := s.storage.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &orderV1.OrderResponse{
		Status: orderV1.OrderStatus_ORDER_STATUS_ACCEPTED,
		Uuid:   req.Uuid,
	}, nil
}

func (s *OrderService) ReturnOrder(ctx context.Context, req *orderV1.ReturnOrderRequest) (*orderV1.OrderResponse, error) {
	order, err := s.storage.FindOne(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	if err := validateReturn(order); err != nil {
		return nil, err
	}

	if err := s.storage.Delete(ctx, req.Uuid); err != nil {
		return nil, err
	}

	return &orderV1.OrderResponse{
		Status: orderV1.OrderStatus_ORDER_STATUS_ACCEPTED,
		Uuid:   req.Uuid,
	}, nil
}

func (s *OrderService) IssueOrder(ctx context.Context, req *orderV1.IssueOrderRequest) (*orderV1.OrderResponse, error) {
	order, err := s.storage.FindOne(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	updOrder := orderV1.Order{
		Uuid:       req.Uuid,
		UserId:     order.UserId,
		Status:     orderV1.OrderStatus_ORDER_STATUS_ISSUED,
		CreatedAt:  order.CreatedAt,
		ExpiresAt:  order.ExpiresAt,
		IssuedAt:   order.IssuedAt,
		ReturnedAt: order.ReturnedAt,
	}

	if err := s.storage.Update(ctx, &updOrder); err != nil {
		return nil, err
	}

	return &orderV1.OrderResponse{
		Status: orderV1.OrderStatus_ORDER_STATUS_ACCEPTED,
		Uuid:   req.Uuid,
	}, nil
}
