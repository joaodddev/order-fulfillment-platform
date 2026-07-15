package query

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Query
type GetOrderStatusQuery struct {
	OrderID string
}

// Query Result
type OrderStatusDTO struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// Query Handler
type GetOrderStatusHandler struct {
	repo order.Repository
}

func NewGetOrderStatusHandler(repo order.Repository) *GetOrderStatusHandler {
	return &GetOrderStatusHandler{repo: repo}
}

func (h *GetOrderStatusHandler) Handle(ctx context.Context, query GetOrderStatusQuery) (*OrderStatusDTO, error) {
	if query.OrderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	orderEntity, err := h.repo.FindByID(ctx, query.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return &OrderStatusDTO{
		OrderID: orderEntity.ID,
		Status:  string(orderEntity.Status),
	}, nil
}
