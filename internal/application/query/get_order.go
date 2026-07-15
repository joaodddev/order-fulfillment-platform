package query

import (
	"context"
	"fmt"
	"time"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Query
type GetOrderQuery struct {
	OrderID string
}

// Query Result (Read Model)
type OrderDTO struct {
	ID          string         `json:"id"`
	CustomerID  string         `json:"customer_id"`
	Status      string         `json:"status"`
	TotalAmount float64        `json:"total_amount"`
	Items       []OrderItemDTO `json:"items"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Query Handler
type GetOrderHandler struct {
	repo order.Repository
}

func NewGetOrderHandler(repo order.Repository) *GetOrderHandler {
	return &GetOrderHandler{repo: repo}
}

func (h *GetOrderHandler) Handle(ctx context.Context, query GetOrderQuery) (*OrderDTO, error) {
	if query.OrderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	orderEntity, err := h.repo.FindByID(ctx, query.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Mapear para DTO
	items := make([]OrderItemDTO, len(orderEntity.Items))
	for i, item := range orderEntity.Items {
		items[i] = OrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &OrderDTO{
		ID:          orderEntity.ID,
		CustomerID:  orderEntity.CustomerID,
		Status:      string(orderEntity.Status),
		TotalAmount: orderEntity.TotalAmount,
		Items:       items,
		CreatedAt:   orderEntity.CreatedAt,
		UpdatedAt:   orderEntity.UpdatedAt,
	}, nil
}
