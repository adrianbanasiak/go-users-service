package users

import (
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/google/uuid"
)

var (
	EvtUserCreated events.EventType = "user:created"

	EvtUserCreatedVersion = 1
)

type EvtUserCreatedPayload struct {
	ID       uuid.UUID
	NickName string
	// ... more properties if needed
}
