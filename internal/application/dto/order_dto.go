package dto

import (
	"time"
)

type CreateOrderRequest struct {
	CustomerID string         `json:"customer_id"`
	Items      []OrderItemDTO `json:"items"`
}

type OrderItemDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderResponse struct {
	ID          string         `json:"id"`
	CustomerID  string         `json:"customer_id"`
	Status      string         `json:"status"`
	TotalAmount float64        `json:"total_amount"`
	Items       []OrderItemDTO `json:"items"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type UpdateOrderRequest struct {
	Items []OrderItemDTO `json:"items"`
}
