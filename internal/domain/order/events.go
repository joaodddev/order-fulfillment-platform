// internal/domain/order/events.go
package order

import (
	"time"
)

// Domain Events for Order aggregate
type OrderCreated struct {
	OrderID    string
	CustomerID string
	Items      []OrderItem
	Total      float64
	Timestamp  time.Time
}

func (e OrderCreated) GetAggregateID() string  { return e.OrderID }
func (e OrderCreated) GetEventType() string    { return "order.created" }
func (e OrderCreated) GetTimestamp() time.Time { return e.Timestamp }

type OrderConfirmed struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderConfirmed) GetAggregateID() string  { return e.OrderID }
func (e OrderConfirmed) GetEventType() string    { return "order.confirmed" }
func (e OrderConfirmed) GetTimestamp() time.Time { return e.Timestamp }

type OrderShipped struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderShipped) GetAggregateID() string  { return e.OrderID }
func (e OrderShipped) GetEventType() string    { return "order.shipped" }
func (e OrderShipped) GetTimestamp() time.Time { return e.Timestamp }

type OrderDelivered struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderDelivered) GetAggregateID() string  { return e.OrderID }
func (e OrderDelivered) GetEventType() string    { return "order.delivered" }
func (e OrderDelivered) GetTimestamp() time.Time { return e.Timestamp }

type OrderCancelled struct {
	OrderID   string
	Reason    string
	Timestamp time.Time
}

func (e OrderCancelled) GetAggregateID() string  { return e.OrderID }
func (e OrderCancelled) GetEventType() string    { return "order.cancelled" }
func (e OrderCancelled) GetTimestamp() time.Time { return e.Timestamp }
