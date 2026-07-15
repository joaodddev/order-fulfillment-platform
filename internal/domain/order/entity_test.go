package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name       string
		customerID string
		items      []OrderItem
		wantErr    error
	}{
		{
			name:       "valid order",
			customerID: "cust-123",
			items: []OrderItem{
				{ProductID: "prod-1", Quantity: 2, Price: 10.5},
				{ProductID: "prod-2", Quantity: 1, Price: 25.0},
			},
			wantErr: nil,
		},
		{
			name:       "empty customer id",
			customerID: "",
			items: []OrderItem{
				{ProductID: "prod-1", Quantity: 1, Price: 10.0},
			},
			wantErr: ErrInvalidCustomerID,
		},
		{
			name:       "empty items",
			customerID: "cust-123",
			items:      []OrderItem{},
			wantErr:    ErrEmptyOrder,
		},
		{
			name:       "invalid quantity",
			customerID: "cust-123",
			items: []OrderItem{
				{ProductID: "prod-1", Quantity: 0, Price: 10.0},
			},
			wantErr: ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.customerID, tt.items)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.NotEmpty(t, order.ID)
				assert.Equal(t, StatusPending, order.Status)
			}
		})
	}
}

func TestOrder_Confirm(t *testing.T) {
	order, _ := NewOrder("cust-123", []OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})

	err := order.Confirm()
	assert.NoError(t, err)
	assert.Equal(t, StatusConfirmed, order.Status)

	// Second confirm should fail
	err = order.Confirm()
	assert.ErrorIs(t, err, ErrOrderAlreadyConfirmed)
}

func TestOrder_Cancel(t *testing.T) {
	order, _ := NewOrder("cust-123", []OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})

	err := order.Cancel()
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, order.Status)

	// Cannot cancel twice
	err = order.Cancel()
	assert.ErrorIs(t, err, ErrOrderAlreadyCancelled)

	// Delivered order cannot be cancelled
	order2, _ := NewOrder("cust-123", []OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})
	order2.Status = StatusDelivered
	err = order2.Cancel()
	assert.ErrorIs(t, err, ErrOrderAlreadyDelivered)
}