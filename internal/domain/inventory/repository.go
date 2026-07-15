package inventory

import "context"

// Repository Interface (Port)
type Repository interface {
	Save(ctx context.Context, inventory *Inventory) error
	FindByID(ctx context.Context, id string) (*Inventory, error)
	FindByProductID(ctx context.Context, productID string) (*Inventory, error)
	Update(ctx context.Context, inventory *Inventory) error
	Delete(ctx context.Context, id string) error
	FindByStatus(ctx context.Context, status StockStatus) ([]*Inventory, error)
}