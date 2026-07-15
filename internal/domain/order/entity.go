package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// Aggregate Root
type Order struct {
	ID          string
	CustomerID  string
	Items       []OrderItem
	Status      OrderStatus
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Value Object
type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

// Value Object
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

// Domain Errors
var (
	ErrInvalidCustomerID = errors.New("customer ID is required")
	ErrEmptyOrder        = errors.New("order must have at least one item")
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
	ErrInvalidPrice      = errors.New("price must be greater than zero")
	ErrOrderAlreadyConfirmed = errors.New("order is already confirmed")
	ErrOrderAlreadyDelivered = errors.New("order is already delivered")
	ErrOrderAlreadyCancelled = errors.New("order is already cancelled")
)

// Factory Method
func NewOrder(customerID string, items []OrderItem) (*Order, error) {
	if customerID == "" {
		return nil, ErrInvalidCustomerID
	}

	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	total := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if item.Price <= 0 {
			return nil, ErrInvalidPrice
		}
		total += item.Price * float64(item.Quantity)
	}

	return &Order{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		Items:       items,
		Status:      StatusPending,
		TotalAmount: total,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Domain Behaviors
func (o *Order) Confirm() error {
	if o.Status != StatusPending {
		return ErrOrderAlreadyConfirmed
	}
	o.Status = StatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Ship() error {
	if o.Status != StatusConfirmed {
		return errors.New("only confirmed orders can be shipped")
	}
	o.Status = StatusShipped
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Deliver() error {
	if o.Status != StatusShipped {
		return errors.New("only shipped orders can be delivered")
	}
	o.Status = StatusDelivered
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusDelivered {
		return ErrOrderAlreadyDelivered
	}
	if o.Status == StatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) AddItem(productID string, quantity int, price float64) error {
	if o.Status != StatusPending {
		return errors.New("cannot modify order after confirmation")
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if price <= 0 {
		return ErrInvalidPrice
	}

	o.Items = append(o.Items, OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	})
	o.TotalAmount += price * float64(quantity)
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) RemoveItem(index int) error {
	if o.Status != StatusPending {
		return errors.New("cannot modify order after confirmation")
	}
	if index < 0 || index >= len(o.Items) {
		return errors.New("invalid item index")
	}

	removed := o.Items[index]
	o.Items = append(o.Items[:index], o.Items[index+1:]...)
	o.TotalAmount -= removed.Price * float64(removed.Quantity)
	o.UpdatedAt = time.Now()
	return nil
}