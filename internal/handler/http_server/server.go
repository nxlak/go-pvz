package openapi

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/nxlak/go-pvz/internal/domain/codes"
	order "github.com/nxlak/go-pvz/internal/repository/storage"
	"github.com/nxlak/go-pvz/pkg/errs"
	order_v1 "github.com/nxlak/go-pvz/pkg/openapi/order/v1"
)

type OrderHandler struct {
	storage order.Repository
}

func NewOrderHandler(storage order.Repository) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) GetOrderById(_ context.Context, params order_v1.GetOrderByIdParams) (order_v1.GetOrderByIdRes, error) {
	order, err := h.storage.FindOne(context.TODO(), params.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (h *OrderHandler) AcceptOrder(_ context.Context, req *order_v1.UpdateOrderRequest, params order_v1.AcceptOrderParams) (order_v1.AcceptOrderRes, error) {
	if err := validateAccept(params.ID, req.UserID, req.ExpiresAt); err != nil {
		return nil, err
	}

	order := order_v1.Order{
		ID:        params.ID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		ExpiresAt: order_v1.NewOptNilDateTime(req.ExpiresAt),
		Status:    order_v1.OrderStatusACCEPTED,
	}

	if err := h.storage.Create(context.TODO(), &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (h *OrderHandler) ReturnOrder(_ context.Context, params order_v1.ReturnOrderParams) (order_v1.ReturnOrderRes, error) {
	order, err := h.storage.FindOne(context.TODO(), params.ID)
	if err != nil {
		return nil, err
	}

	if err := validateReturn(order); err != nil {
		return nil, err
	}

	if err := h.storage.Delete(context.TODO(), params.ID); err != nil {
		return nil, err
	}

	return &order_v1.ReturnOrderNoContent{}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *order_v1.PatchOrderRequest, params order_v1.UpdateOrderParams) (order_v1.UpdateOrderRes, error) {
	order, err := h.storage.FindOne(context.TODO(), params.ID)
	if err != nil {
		return nil, err
	}

	updOrder := order_v1.Order{
		ID:         params.ID,
		UserID:     order.UserID,
		Status:     order_v1.OrderStatus(req.Status),
		CreatedAt:  order.CreatedAt,
		ExpiresAt:  order.ExpiresAt,
		IssuedAt:   order.IssuedAt,
		ReturnedAt: order.ReturnedAt,
	}

	if err := h.storage.Update(context.TODO(), &updOrder); err != nil {
		return nil, err
	}

	return &updOrder, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *order_v1.AppErrorStatusCode {
	status := http.StatusInternalServerError

	apiErr := order_v1.AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal error",
	}

	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		apiErr.Code = appErr.Code
		apiErr.Message = appErr.Message
		apiErr.Fields = errs.ToAPIErrorFields(appErr.Fields)

		switch appErr.Code {
		case codes.ErrOrderNotFound.Code:
			status = http.StatusNotFound
		case codes.ErrOrderAlreadyExists.Code:
			status = http.StatusConflict
		case codes.ErrValidationFailed.Code:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
	}

	return &order_v1.AppErrorStatusCode{
		StatusCode: status,
		Response:   apiErr,
	}
}
