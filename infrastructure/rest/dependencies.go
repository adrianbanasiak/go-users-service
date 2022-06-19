package rest

import (
	"github.com/adrianbanasiak/go-users-service/infrastructure/healthchecks"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/adrianbanasiak/go-users-service/internal/users"
)

type Dependencies struct {
	Log                 shared.Logger
	UsersService        *users.Service
	HealthchecksService *healthchecks.Service
}
