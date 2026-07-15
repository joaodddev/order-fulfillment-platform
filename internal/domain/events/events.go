package events

import (
	"time"
)

// Base Event
type DomainEvent interface {
	GetAggregateID() string
	GetEventType() string
	GetTimestamp() time.Time
}

// Order Events
type OrderCreatedEvent struct {
	OrderID    string
	CustomerID string
	Items      []OrderItemEvent
	Total      float64
	Timestamp  time.Time
}

type OrderItemEvent struct {
	ProductID string
	Quantity  int
	Price     float64
}

func (e OrderCreatedEvent) GetAggregateID() string { return e.OrderID }
func (e OrderCreatedEvent) GetEventType() string   { return "order.created" }
func (e OrderCreatedEvent) GetTimestamp() time.Time { return e.Timestamp }

type OrderConfirmedEvent struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderConfirmedEvent) GetAggregateID() string { return e.OrderID }
func (e OrderConfirmedEvent) GetEventType() string   { return "order.confirmed" }
func (e OrderConfirmedEvent) GetTimestamp() time.Time { return e.Timestamp }

type OrderShippedEvent struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderShippedEvent) GetAggregateID() string { return e.OrderID }
func (e OrderShippedEvent) GetEventType() string   { return "order.shipped" }
func (e OrderShippedEvent) GetTimestamp() time.Time { return e.Timestamp }

type OrderDeliveredEvent struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderDeliveredEvent) GetAggregateID() string { return e.OrderID }
func (e OrderDeliveredEvent) GetEventType() string   { return "order.delivered" }
func (e OrderDeliveredEvent) GetTimestamp() time.Time { return e.Timestamp }

type OrderCancelledEvent struct {
	OrderID   string
	Reason    string
	Timestamp time.Time
}

func (e OrderCancelledEvent) GetAggregateID() string { return e.OrderID }
func (e OrderCancelledEvent) GetEventType() string   { return "order.cancelled" }
func (e OrderCancelledEvent) GetTimestamp() time.Time { return e.Timestamp }

// Inventory Events
type StockReservedEvent struct {
	InventoryID string
	ProductID   string
	Quantity    int
	Timestamp   time.Time
}

func (e StockReservedEvent) GetAggregateID() string { return e.InventoryID }
func (e StockReservedEvent) GetEventType() string   { return "stock.reserved" }
func (e StockReservedEvent) GetTimestamp() time.Time { return e.Timestamp }

type StockReleasedEvent struct {
	InventoryID string
	ProductID   string
	Quantity    int
	Timestamp   time.Time
}

func (e StockReleasedEvent) GetAggregateID() string { return e.InventoryID }
func (e StockReleasedEvent) GetEventType() string   { return "stock.released" }
func (e StockReleasedEvent) GetTimestamp() time.Time { return e.Timestamp }

type StockUpdatedEvent struct {
	InventoryID string
	ProductID   string
	OldQuantity int
	NewQuantity int
	Timestamp   time.Time
}

func (e StockUpdatedEvent) GetAggregateID() string { return e.InventoryID }
func (e StockUpdatedEvent) GetEventType() string   { return "stock.updated" }
func (e StockUpdatedEvent) GetTimestamp() time.Time { return e.Timestamp }