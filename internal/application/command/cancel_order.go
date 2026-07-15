package command

import (
	"context"
	"fmt"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/joaodddev/order-fulfillment-platform/internal/domain/order"
)

// Command
type CancelOrderCommand struct {
	OrderID string
	Reason  string
}

// Command Handler
type CancelOrderHandler struct {
	orderRepo      order.Repository
	eventPublisher events.EventPublisher
}

func NewCancelOrderHandler(
	orderRepo order.Repository,
	eventPublisher events.EventPublisher,
) *CancelOrderHandler {
	return &CancelOrderHandler{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (h *CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrderCommand) error {
	// 1. Validar
	if cmd.OrderID == "" {
		return fmt.Errorf("order ID is required")
	}

	// 2. Buscar agregado
	existingOrder, err := h.orderRepo.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 3. Executar comportamento de domínio
	if err := existingOrder.Cancel(); err != nil {
		return err
	}

	// 4. Persistir
	if err := h.orderRepo.Update(ctx, existingOrder); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// 5. Publicar evento
	event := order.OrderCancelled{
		OrderID:   existingOrder.ID,
		Reason:    cmd.Reason,
		Timestamp: existingOrder.UpdatedAt,
	}
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		fmt.Printf("failed to publish event: %v\n", err)
	}

	return nil
}
