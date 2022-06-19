package users

import (
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/google/uuid"
)

var (
	EvtUserCreated      events.EventType = "user:created"
	EvtUserDeleted      events.EventType = "user:deleted"
	EvtUserEmailUpdated events.EventType = "user:email_updated"

	EvtUserCreatedVersion      = 1
	EvtUserDeletedVersion      = 1
	EvtUserEmailUpdatedVersion = 1
)

type EvtUserCreatedPayload struct {
	ID       uuid.UUID
	NickName string
	// ... more properties if needed
}

type EvtUserDeletedPayload struct {
	ID uuid.UUID
}

type EvtUserEmailUpdatedPayload struct {
	ID  uuid.UUID
	Old string
	New string
}
