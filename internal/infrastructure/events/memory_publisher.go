package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
)

// In-memory event publisher (for development)
type MemoryEventPublisher struct {
	subscribers map[string][]events.EventHandler
	mu          sync.RWMutex
}

func NewMemoryEventPublisher() *MemoryEventPublisher {
	return &MemoryEventPublisher{
		subscribers: make(map[string][]events.EventHandler),
	}
}

func (p *MemoryEventPublisher) Publish(ctx context.Context, event events.DomainEvent) error {
	eventType := event.GetEventType()
	p.mu.RLock()
	handlers, exists := p.subscribers[eventType]
	p.mu.RUnlock()

	if !exists {
		return nil // No subscribers, just return
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return fmt.Errorf("error handling event %s: %w", eventType, err)
		}
	}

	return nil
}

func (p *MemoryEventPublisher) PublishBatch(ctx context.Context, events []events.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (p *MemoryEventPublisher) Subscribe(eventType string, handler events.EventHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscribers[eventType] = append(p.subscribers[eventType], handler)
	return nil
}

// Convenience function for subscribing with function
func (p *MemoryEventPublisher) SubscribeFunc(eventType string, fn func(ctx context.Context, event events.DomainEvent) error) error {
	return p.Subscribe(eventType, events.EventHandlerFunc(fn))
}