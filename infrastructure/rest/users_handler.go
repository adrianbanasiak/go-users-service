package rest

import (
	"encoding/json"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/adrianbanasiak/go-users-service/internal/users"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func NewUsersHandler(log shared.Logger, router *mux.Router, service *users.Service) *UsersHandler {
	h := UsersHandler{
		log:     log,
		service: service,
		router:  router,
		prefix:  "/api/users",
	}

	h.RegisterEndpoints()

	return &h
}

type UsersHandler struct {
	prefix  string
	router  *mux.Router
	log     shared.Logger
	service *users.Service
}

func (h *UsersHandler) RegisterEndpoints() {
	h.router.HandleFunc(h.prefix, h.handleCreate).
		Methods(http.MethodPost)

	h.log.Infow("registered endpoint", "endpoint", "users")
}

func (h *UsersHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.sendError(w, err)
		return
	}

	// could validate request data with go-validator
	var req CreateUserReq
	err = json.Unmarshal(b, &req)
	if err != nil {
		h.sendError(w, err)
		return
	}

	user, err := h.service.CreateUser(r.Context(), users.CreateUserReq{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		NickName:    req.NickName,
		Password:    req.Password,
		Email:       req.Email,
		CountryCode: req.CountryCode,
	})

	if err != nil {
		h.sendError(w, err)
		return
	}

	err = RespondSuccess(w, PresentUser(user), http.StatusCreated)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users:create",
			"error", err)
		return
	}
}

func (h *UsersHandler) sendError(w http.ResponseWriter, err error) {
	// should have some kind of error mapper to vary status code by type of error
	err = RespondError(w, []error{err}, http.StatusPreconditionFailed)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users:create",
			"error", err)
	}
}
