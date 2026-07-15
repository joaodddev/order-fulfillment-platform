package command

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Command
type DeliverOrderCommand struct {
	OrderID string
}

// Command Handler
type DeliverOrderHandler struct {
	orderRepo      order.Repository
	eventPublisher events.EventPublisher
}

func NewDeliverOrderHandler(
	orderRepo order.Repository,
	eventPublisher events.EventPublisher,
) *DeliverOrderHandler {
	return &DeliverOrderHandler{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (h *DeliverOrderHandler) Handle(ctx context.Context, cmd DeliverOrderCommand) error {
	if cmd.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}

	existingOrder, err := h.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if err := existingOrder.Deliver(); err != nil {
		return err
	}

	if err := h.orderRepo.Update(ctx, existingOrder); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	event := order.OrderDelivered{
		OrderID:   existingOrder.ID,
		Timestamp: existingOrder.UpdatedAt,
	}
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish event: %v\n", err)
	}

	return nil
}
