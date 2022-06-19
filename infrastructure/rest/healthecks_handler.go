package rest

import (
	"github.com/adrianbanasiak/go-users-service/infrastructure/healthchecks"
	"github.com/gorilla/mux"
)

func NewHealthchecksHandler(service *healthchecks.Service, router *mux.Router) *HealthchecksHandler {
	h := HealthchecksHandler{service: service, router: router, prefix: "/health"}

	h.Register()
	return &h
}

type HealthchecksHandler struct {
	service *healthchecks.Service
	router  *mux.Router
	prefix  string
}

func (h *HealthchecksHandler) Register() {
	h.router.HandleFunc(h.prefix, h.service.ConfigureHandler())
}
