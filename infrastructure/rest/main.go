package rest

import (
	"fmt"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	listenPort   int
	router       *mux.Router
	log          shared.Logger
	dependencies Dependencies
}

func NewServer(deps Dependencies, listenPort int) *Server {
	s := Server{
		listenPort:   listenPort,
		log:          deps.Log,
		router:       mux.NewRouter(),
		dependencies: deps,
	}

	s.setupHandlers()

	return &s
}

func (s *Server) Start() error {
	s.log.Info("starting REST API server")

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", s.listenPort),
		Handler: s.router,
	}

	return srv.ListenAndServe()
}

func (s *Server) setupHandlers() {
	_ = NewUsersHandler(s.log, s.router, s.dependencies.UsersService)
	_ = NewHealthchecksHandler(s.dependencies.HealthchecksService, s.router)
}
