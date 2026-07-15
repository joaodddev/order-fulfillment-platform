package events

import (
	"context"
)

// Event Publisher Interface (Port)
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, events []DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
}

// Event Handler Interface
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// Event Handler Function (convenience)
type EventHandlerFunc func(ctx context.Context, event DomainEvent) error

func (f EventHandlerFunc) Handle(ctx context.Context, event DomainEvent) error {
	return f(ctx, event)
}
