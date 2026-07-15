package order

import "context"

// Repository Interface (Port) - following DIP
type Repository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
	FindByCustomerID(ctx context.Context, customerID string) ([]*Order, error)
	FindByStatus(ctx context.Context, status OrderStatus) ([]*Order, error)
}