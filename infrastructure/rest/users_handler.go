package rest

import (
	"encoding/json"
	"errors"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/adrianbanasiak/go-users-service/internal/users"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
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
	h.router.HandleFunc(h.prefix, h.handleList).
		Queries("page", "{page:[0-9]+}", "size", "{size:[0-9]+}").
		Methods(http.MethodGet)
	h.router.HandleFunc(h.prefix+"/{id}", h.handleDelete).
		Methods(http.MethodDelete)

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

func (h *UsersHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]

	if !ok {
		h.sendError(w, errors.New("missing ID path parameter"))
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		h.sendError(w, errors.New("invalid ID path parameter"))
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		h.sendError(w, err)
		return
	}

	_ = RespondSuccess(w, nil, http.StatusNoContent)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users:delete",
			"error", err)
		return
	}
}

func (h *UsersHandler) handleList(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.sendError(w, errors.New("invalid query string param page"))
		return
	}

	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		h.sendError(w, errors.New("invalid query string param size"))
		return
	}

	uu, err := h.service.ListUsers(r.Context(), users.ListUsersReq{
		Page: page,
		Size: size,
	})

	if err != nil {
		h.sendError(w, err)
	}

	res := make([]UserRes, 0)
	for _, user := range uu {
		res = append(res, PresentUser(user))
	}

	err = RespondSuccess(w, res, http.StatusOK)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "uu:list",
			"error", err)
		return
	}
}

func (h *UsersHandler) sendError(w http.ResponseWriter, err error) {
	// should have some kind of error mapper to vary status code by type of error
	err = RespondError(w, []error{err}, http.StatusPreconditionFailed)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users",
			"error", err)
	}
}
