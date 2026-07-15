package query

import (
	"context"
	"testing"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Order Repository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(ctx context.Context, order *order.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *order.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*order.Order, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status order.OrderStatus) ([]*order.Order, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func TestGetOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	handler := NewGetOrderHandler(mockRepo)

	// Criar um pedido de teste
	testOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 2, Price: 10.0},
	})

	// Configurar mock
	mockRepo.On("FindByID", mock.Anything, testOrder.ID).Return(testOrder, nil)

	// Executar query
	query := GetOrderQuery{OrderID: testOrder.ID}
	result, err := handler.Handle(context.Background(), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testOrder.ID, result.ID)
	assert.Equal(t, string(testOrder.Status), result.Status)
	assert.Equal(t, testOrder.TotalAmount, result.TotalAmount)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "prod-1", result.Items[0].ProductID)

	mockRepo.AssertExpectations(t)
}

func TestGetOrderHandler_NotFound(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	handler := NewGetOrderHandler(mockRepo)

	mockRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, assert.AnError)

	query := GetOrderQuery{OrderID: "non-existent"}
	result, err := handler.Handle(context.Background(), query)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")
}

func TestListOrdersHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	handler := NewListOrdersHandler(mockRepo)

	// Criar pedidos de teste
	order1, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})
	order2, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-2", Quantity: 2, Price: 20.0},
	})

	orders := []*order.Order{order1, order2}

	// Configurar mock
	mockRepo.On("FindByCustomerID", mock.Anything, "cust-123").Return(orders, nil)

	// Executar query
	query := ListOrdersQuery{
		CustomerID: "cust-123",
		Limit:      10,
		Offset:     0,
	}
	result, err := handler.Handle(context.Background(), query)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, order1.ID, result[0].ID)
	assert.Equal(t, order2.ID, result[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestGetOrderStatusHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	handler := NewGetOrderStatusHandler(mockRepo)

	testOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})

	mockRepo.On("FindByID", mock.Anything, testOrder.ID).Return(testOrder, nil)

	query := GetOrderStatusQuery{OrderID: testOrder.ID}
	result, err := handler.Handle(context.Background(), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testOrder.ID, result.OrderID)
	assert.Equal(t, string(testOrder.Status), result.Status)

	mockRepo.AssertExpectations(t)
}
