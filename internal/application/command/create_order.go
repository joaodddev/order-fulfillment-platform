package command

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Command
type CreateOrderCommand struct {
	CustomerID string
	Items      []OrderItemCommand
}

type OrderItemCommand struct {
	ProductID string
	Quantity  int
	Price     float64
}

// Command Handler
type CreateOrderHandler struct {
	orderRepo      order.Repository
	eventPublisher events.EventPublisher
}

func NewCreateOrderHandler(
	orderRepo order.Repository,
	eventPublisher events.EventPublisher,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (string, error) {
	// 1. Validar comando
	if cmd.CustomerID == "" {
		return "", fmt.Errorf("customer ID is required")
	}
	if len(cmd.Items) == 0 {
		return "", fmt.Errorf("order must have at least one item")
	}

	// 2. Criar itens do domínio
	items := make([]order.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		items[i] = order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// 3. Criar agregado (Domain Factory)
	newOrder, err := order.NewOrder(cmd.CustomerID, items)
	if err != nil {
		return "", err
	}

	// 4. Persistir
	if err := h.orderRepo.Save(ctx, newOrder); err != nil {
		return "", fmt.Errorf("failed to save order: %w", err)
	}

	// 5. Publicar evento de domínio
	event := order.OrderCreated{
		OrderID:    newOrder.ID,
		CustomerID: newOrder.CustomerID,
		Items:      newOrder.Items,
		Total:      newOrder.TotalAmount,
		Timestamp:  newOrder.CreatedAt,
	}
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		// In production, you might want to use a proper logger
		fmt.Printf("failed to publish event: %v\n", err)
	}

	return newOrder.ID, nil
}
