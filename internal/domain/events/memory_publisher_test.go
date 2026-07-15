package events

import (
	"context"
	"testing"
	"time"

	"github.com/joaodddev/order-fulfillment-platform/internal/domain/events"
	"github.com/stretchr/testify/assert"
)

func TestMemoryEventPublisher(t *testing.T) {
	publisher := NewMemoryEventPublisher()

	// Test subscription
	eventReceived := false
	err := publisher.SubscribeFunc("test.event", func(ctx context.Context, event events.DomainEvent) error {
		eventReceived = true
		return nil
	})
	assert.NoError(t, err)

	// Test publishing
	testEvent := struct {
		events.DomainEvent
		ID string
	}{
		ID: "test-123",
	}
	// Wrapping as DomainEvent
	type TestEvent struct {
		ID        string
		Timestamp time.Time
	}
	testEventWrapper := struct {
		events.DomainEvent
		TestEvent
	}{
		TestEvent: TestEvent{
			ID:        "test-123",
			Timestamp: time.Now(),
		},
	}

	// For testing, we'll use a simpler approach
	// We need to implement GetAggregateID, GetEventType, GetTimestamp
	// Let's create a real test event

	type SimpleEvent struct {
		AggregateID string
		EventType   string
		Timestamp   time.Time
	}

	func (e SimpleEvent) GetAggregateID() string { return e.AggregateID }
	func (e SimpleEvent) GetEventType() string   { return e.EventType }
	func (e SimpleEvent) GetTimestamp() time.Time { return e.Timestamp }

	event := SimpleEvent{
		AggregateID: "agg-123",
		EventType:   "test.event",
		Timestamp:   time.Now(),
	}

	err = publisher.Publish(context.Background(), event)
	assert.NoError(t, err)
	assert.True(t, eventReceived)
}

func TestMemoryEventPublisher_MultipleSubscribers(t *testing.T) {
	publisher := NewMemoryEventPublisher()

	counter := 0
	for i := 0; i < 3; i++ {
		err := publisher.SubscribeFunc("test.event", func(ctx context.Context, event events.DomainEvent) error {
			counter++
			return nil
		})
		assert.NoError(t, err)
	}

	type SimpleEvent struct {
		AggregateID string
		EventType   string
		Timestamp   time.Time
	}

	func (e SimpleEvent) GetAggregateID() string { return e.AggregateID }
	func (e SimpleEvent) GetEventType() string   { return e.EventType }
	func (e SimpleEvent) GetTimestamp() time.Time { return e.Timestamp }

	event := SimpleEvent{
		AggregateID: "agg-123",
		EventType:   "test.event",
		Timestamp:   time.Now(),
	}

	err := publisher.Publish(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, 3, counter)
}