package users

import (
	"context"
	"errors"
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
)

var (
	ErrCreateFailed = errors.New("failed to create user")
)

func NewService(r Repository, log shared.Logger, bus events.Bus) *Service {
	return &Service{repository: r, log: log, bus: bus}
}

type Service struct {
	repository Repository
	log        shared.Logger
	bus        events.Bus
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserReq) (User, error) {
	s.log.Infow("create user", "nickName", req.NickName)

	u, err := NewUser(req)
	if err != nil {
		s.log.Errorw("failed to create user",
			"error", err)
		return User{}, ErrCreateFailed
	}

	created, err := s.repository.Create(ctx, u)
	if err != nil {
		s.log.Errorw("failed to persist user through repository",
			"error", err)
		return User{}, ErrCreateFailed
	}

	evt := events.NewEvent(EvtUserCreated, EvtUserCreatedVersion,
		EvtUserCreatedPayload{ID: created.ID, NickName: created.NickName})

	if err = s.bus.Enqueue(evt); err != nil {
		s.log.Errorw("failed to emit user created event",
			"error", err,
			"userID", created.ID)
		return User{}, ErrCreateFailed
	}

	s.log.Infow("user created successfully",
		"userID", created.ID,
		"nickName", created.NickName)

	return created, err
}
