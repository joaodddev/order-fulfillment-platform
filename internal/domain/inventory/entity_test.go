package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInventory(t *testing.T) {
	tests := []struct {
		name      string
		productID string
		quantity  int
		minStock  int
		maxStock  int
		wantErr   error
	}{
		{
			name:      "valid inventory",
			productID: "prod-123",
			quantity:  100,
			minStock:  10,
			maxStock:  200,
			wantErr:   nil,
		},
		{
			name:      "empty product id",
			productID: "",
			quantity:  100,
			minStock:  10,
			maxStock:  200,
			wantErr:   errors.New("product ID is required"),
		},
		{
			name:      "negative quantity",
			productID: "prod-123",
			quantity:  -5,
			minStock:  10,
			maxStock:  200,
			wantErr:   ErrInvalidQuantity,
		},
		{
			name:      "invalid min stock",
			productID: "prod-123",
			quantity:  100,
			minStock:  -5,
			maxStock:  200,
			wantErr:   ErrInvalidMinStock,
		},
		{
			name:      "invalid max stock",
			productID: "prod-123",
			quantity:  100,
			minStock:  10,
			maxStock:  5,
			wantErr:   ErrInvalidMaxStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, err := NewInventory(tt.productID, tt.quantity, tt.minStock, tt.maxStock)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, inv)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, inv)
				assert.NotEmpty(t, inv.ID)
				assert.Equal(t, tt.quantity, inv.Quantity)
				assert.Equal(t, 0, inv.Reserved)
				assert.Equal(t, StatusInStock, inv.Status)
			}
		})
	}
}

func TestInventory_ReserveStock(t *testing.T) {
	inv, _ := NewInventory("prod-123", 100, 10, 200)

	// Test successful reservation
	err := inv.ReserveStock(30)
	assert.NoError(t, err)
	assert.Equal(t, 30, inv.Reserved)
	assert.Equal(t, 70, inv.GetAvailableStock())

	// Test insufficient stock
	err = inv.ReserveStock(80)
	assert.ErrorIs(t, err, ErrInsufficientStock)

	// Test invalid quantity
	err = inv.ReserveStock(-5)
	assert.ErrorIs(t, err, ErrInvalidQuantity)
}

func TestInventory_ReleaseStock(t *testing.T) {
	inv, _ := NewInventory("prod-123", 100, 10, 200)
	inv.ReserveStock(30)

	// Test successful release
	err := inv.ReleaseStock(20)
	assert.NoError(t, err)
	assert.Equal(t, 10, inv.Reserved)
	assert.Equal(t, 90, inv.GetAvailableStock())

	// Test releasing more than reserved
	err = inv.ReleaseStock(20)
	assert.Error(t, err)

	// Test invalid quantity
	err = inv.ReleaseStock(-5)
	assert.ErrorIs(t, err, ErrInvalidQuantity)
}

func TestInventory_StatusUpdate(t *testing.T) {
	inv, _ := NewInventory("prod-123", 5, 10, 200)
	assert.Equal(t, StatusLowStock, inv.Status)

	inv.AddStock(10)
	assert.Equal(t, StatusInStock, inv.Status)

	inv.RemoveStock(20)
	assert.Equal(t, StatusOutOfStock, inv.Status)
}