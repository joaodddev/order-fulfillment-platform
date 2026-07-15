package command

import (
	"context"
	"testing"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
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

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(ctx context.Context, event events.DomainEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishBatch(ctx context.Context, events []events.DomainEvent) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func (m *MockEventPublisher) Subscribe(eventType string, handler events.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func TestCreateOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	handler := NewCreateOrderHandler(mockRepo, mockPublisher)

	cmd := CreateOrderCommand{
		CustomerID: "cust-123",
		Items: []OrderItemCommand{
			{ProductID: "prod-1", Quantity: 2, Price: 10.0},
		},
	}

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*order.Order")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("order.OrderCreated")).Return(nil)

	orderID, err := handler.Handle(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotEmpty(t, orderID)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestConfirmOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	existingOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})

	handler := NewConfirmOrderHandler(mockRepo, mockPublisher)

	mockRepo.On("FindByID", mock.Anything, existingOrder.ID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*order.Order")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("order.OrderConfirmed")).Return(nil)

	err := handler.Handle(context.Background(), ConfirmOrderCommand{OrderID: existingOrder.ID})

	assert.NoError(t, err)
	assert.Equal(t, order.StatusConfirmed, existingOrder.Status)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestCancelOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	existingOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})

	handler := NewCancelOrderHandler(mockRepo, mockPublisher)

	mockRepo.On("FindByID", mock.Anything, existingOrder.ID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*order.Order")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("order.OrderCancelled")).Return(nil)

	err := handler.Handle(context.Background(), CancelOrderCommand{
		OrderID: existingOrder.ID,
		Reason:  "Customer requested",
	})

	assert.NoError(t, err)
	assert.Equal(t, order.StatusCancelled, existingOrder.Status)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestShipOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	existingOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})
	existingOrder.Confirm() // Confirm first

	handler := NewShipOrderHandler(mockRepo, mockPublisher)

	mockRepo.On("FindByID", mock.Anything, existingOrder.ID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*order.Order")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("order.OrderShipped")).Return(nil)

	err := handler.Handle(context.Background(), ShipOrderCommand{OrderID: existingOrder.ID})

	assert.NoError(t, err)
	assert.Equal(t, order.StatusShipped, existingOrder.Status)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestDeliverOrderHandler(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	existingOrder, _ := order.NewOrder("cust-123", []order.OrderItem{
		{ProductID: "prod-1", Quantity: 1, Price: 10.0},
	})
	existingOrder.Confirm()
	existingOrder.Ship()

	handler := NewDeliverOrderHandler(mockRepo, mockPublisher)

	mockRepo.On("FindByID", mock.Anything, existingOrder.ID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*order.Order")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("order.OrderDelivered")).Return(nil)

	err := handler.Handle(context.Background(), DeliverOrderCommand{OrderID: existingOrder.ID})

	assert.NoError(t, err)
	assert.Equal(t, order.StatusDelivered, existingOrder.Status)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
