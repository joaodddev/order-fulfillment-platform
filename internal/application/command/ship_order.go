package command

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Command
type ShipOrderCommand struct {
	OrderID string
}

// Command Handler
type ShipOrderHandler struct {
	orderRepo      order.Repository
	eventPublisher events.EventPublisher
}

func NewShipOrderHandler(
	orderRepo order.Repository,
	eventPublisher events.EventPublisher,
) *ShipOrderHandler {
	return &ShipOrderHandler{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (h *ShipOrderHandler) Handle(ctx context.Context, cmd ShipOrderCommand) error {
	if cmd.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}

	existingOrder, err := h.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if err := existingOrder.Ship(); err != nil {
		return err
	}

	if err := h.orderRepo.Update(ctx, existingOrder); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	event := order.OrderShipped{
		OrderID:   existingOrder.ID,
		Timestamp: existingOrder.UpdatedAt,
	}
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish event: %v\n", err)
	}

	return nil
}
