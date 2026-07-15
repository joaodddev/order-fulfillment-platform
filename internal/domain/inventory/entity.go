package inventory

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type StockStatus string

const (
	StatusInStock     StockStatus = "in_stock"
	StatusLowStock    StockStatus = "low_stock"
	StatusOutOfStock  StockStatus = "out_of_stock"
)

// Aggregate Root
type Inventory struct {
	ID          string
	ProductID   string
	Quantity    int
	Reserved    int
	MinStock    int
	MaxStock    int
	Status      StockStatus
	LastUpdated time.Time
}

// Value Object
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	SKU         string
}

// Domain Errors
var (
	ErrProductNotFound       = errors.New("product not found in inventory")
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrInvalidQuantity       = errors.New("quantity must be greater than zero")
	ErrInvalidMinStock       = errors.New("minimum stock must be greater than zero")
	ErrInvalidMaxStock       = errors.New("maximum stock must be greater than minimum stock")
	ErrStockAlreadyReserved  = errors.New("stock already reserved")
)

// Factory
func NewInventory(productID string, quantity, minStock, maxStock int) (*Inventory, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if minStock < 0 {
		return nil, ErrInvalidMinStock
	}
	if maxStock <= minStock {
		return nil, ErrInvalidMaxStock
	}

	inv := &Inventory{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    quantity,
		Reserved:    0,
		MinStock:    minStock,
		MaxStock:    maxStock,
		LastUpdated: time.Now(),
	}
	inv.updateStatus()
	return inv, nil
}

// Domain Behaviors
func (i *Inventory) ReserveStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	available := i.Quantity - i.Reserved
	if available < quantity {
		return ErrInsufficientStock
	}

	i.Reserved += quantity
	i.LastUpdated = time.Now()
	i.updateStatus()
	return nil
}

func (i *Inventory) ReleaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.Reserved < quantity {
		return errors.New("cannot release more than reserved")
	}

	i.Reserved -= quantity
	i.LastUpdated = time.Now()
	i.updateStatus()
	return nil
}

func (i *Inventory) AddStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	i.Quantity += quantity
	i.LastUpdated = time.Now()
	i.updateStatus()
	return nil
}

func (i *Inventory) RemoveStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if i.Quantity-quantity < 0 {
		return ErrInsufficientStock
	}

	i.Quantity -= quantity
	i.LastUpdated = time.Now()
	i.updateStatus()
	return nil
}

func (i *Inventory) CheckAvailability(quantity int) bool {
	if quantity <= 0 {
		return false
	}
	return (i.Quantity - i.Reserved) >= quantity
}

func (i *Inventory) GetAvailableStock() int {
	return i.Quantity - i.Reserved
}

// Private methods
func (i *Inventory) updateStatus() {
	available := i.GetAvailableStock()
	switch {
	case available <= 0:
		i.Status = StatusOutOfStock
	case available <= i.MinStock:
		i.Status = StatusLowStock
	default:
		i.Status = StatusInStock
	}
}