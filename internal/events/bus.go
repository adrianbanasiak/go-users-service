package events

import "github.com/google/uuid"

type Bus interface {
	// Enqueue allows to notify other parts of the system about domain event which occurred in the system.
	// Caller should compensate/rollback operation if enqueue fails
	Enqueue(evt Event) error
}

type EventType string

type Event struct {
	ID      uuid.UUID
	Type    EventType
	Version int
	Payload any
}
