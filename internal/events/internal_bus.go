package events

import (
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/google/uuid"
)

func NewInternalEventBus(log shared.Logger) *InternalEventBus {
	return &InternalEventBus{log: log}
}

type InternalEventBus struct {
	log shared.Logger
}

func NewEvent(eventType EventType, version int, payload any) Event {
	return Event{ID: uuid.New(), Type: eventType, Version: version, Payload: payload}
}

func (s *InternalEventBus) Enqueue(evt Event) error {
	s.log.Infow("enqueued event",
		"eventID", evt.ID,
		"eventType", evt.Type,
		"version", evt.Version)

	return nil
}
