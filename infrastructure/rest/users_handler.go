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

	h.errorsMap = map[error]int{
		users.ErrNotFound: http.StatusNotFound,
	}

	return &h
}

type UsersHandler struct {
	prefix    string
	router    *mux.Router
	log       shared.Logger
	service   *users.Service
	errorsMap map[error]int
}

func (h *UsersHandler) RegisterEndpoints() {
	h.router.HandleFunc(h.prefix, h.handleCreate).
		Methods(http.MethodPost)
	h.router.HandleFunc(h.prefix, h.handleList).
		Queries("page", "{page:[0-9]+}", "size", "{size:[0-9]+}").
		Methods(http.MethodGet)
	h.router.HandleFunc(h.prefix+"/{id}", h.handleDelete).
		Methods(http.MethodDelete)
	h.router.HandleFunc(h.prefix+"/{id}/change-email", h.handleChangeEmail).
		Methods(http.MethodPatch)

	h.log.Infow("registered endpoint", "endpoint", "users")
}

func (h *UsersHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorw("failed to read request body",
			"error", err,
			"handler", "users:create")
		h.sendError(w, err)
		return
	}

	// could validate request data with go-validator
	var req CreateUserReq
	err = json.Unmarshal(b, &req)
	if err != nil {
		h.log.Errorw("failed to decode request body to JSON",
			"error", err,
			"handler", "users:create")
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

	err = RespondSuccess(w, nil, http.StatusNoContent)
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

func (h *UsersHandler) handleChangeEmail(w http.ResponseWriter, r *http.Request) {
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

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorw("failed to read request body",
			"error", err,
			"handler", "users:change_email")
		h.sendError(w, err)
		return
	}

	var req UserChangeEmailReq
	err = json.Unmarshal(b, &req)
	if err != nil {
		h.log.Errorw("failed to decode request body to JSON",
			"error", err,
			"handler", "users:change_email")
		h.sendError(w, err)
		return
	}

	u, err := h.service.ChangeEmail(r.Context(), userID, req.Email)
	if err != nil {
		h.sendError(w, err)
		return
	}

	err = RespondSuccess(w, PresentUser(u), http.StatusOK)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users:change_email",
			"error", err)
		return
	}
}

func (h *UsersHandler) sendError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError

	// simple error mapper
	if sc, ok := h.errorsMap[err]; ok {
		statusCode = sc
	}

	err = RespondError(w, []error{err}, statusCode)
	if err != nil {
		h.log.Errorw("failed to write response to the client",
			"handler", "users",
			"error", err)
	}
}
