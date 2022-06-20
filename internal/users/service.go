package users

import (
	"context"
	"errors"
	"github.com/adrianbanasiak/go-users-service/internal/events"
	"github.com/adrianbanasiak/go-users-service/internal/shared"
	"github.com/google/uuid"
)

var (
	ErrCreateFailed           = errors.New("failed to create user")
	ErrEmailInUse             = errors.New("email already used")
	ErrActionFailed           = errors.New("action has failed")
	ErrPaginationSizeTooBig   = errors.New("pagination size is too big")
	ErrInvalidPaginationRange = errors.New("invalid pagination range")
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

	_, err := s.repository.FindByEmail(ctx, req.Email)
	switch err {
	case ErrNotFound:
	case nil:
		s.log.Errorw("failed to create user",
			"error", ErrEmailInUse)
		return User{}, ErrEmailInUse
	default:
		return User{}, err
	}

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

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	s.log.Infow("delete user", "userID", userID)

	err := s.repository.Delete(ctx, userID)
	switch err {
	case nil:
	case ErrNotFound:
		s.log.Infow("failed to delete user - it was already removed or does not exist")
		return nil
	default:
		s.log.Errorw("failed to delete user",
			"userID", userID)
		return err

	}

	evt := events.NewEvent(EvtUserDeleted, EvtUserDeletedVersion, EvtUserDeletedPayload{ID: userID})
	err = s.bus.Enqueue(evt)
	if err != nil {
		s.log.Errorw("failed to emit user deleted event",
			"error", err,
			"userID", userID)
		return ErrActionFailed
	}

	s.log.Infow("user deleted successfully",
		"userID", userID)
	return nil
}

func (s *Service) ListUsers(ctx context.Context, req ListUsersReq) ([]User, error) {
	s.log.Infow("list users with pagination")

	if req.Size < 1 || req.Page < 1 {
		s.log.Errorw("failed to list users",
			"error", ErrInvalidPaginationRange)
		return nil, ErrInvalidPaginationRange
	}

	if req.Size > 100 {
		s.log.Errorw("failed to list users",
			"error", ErrPaginationSizeTooBig)
		return nil, ErrPaginationSizeTooBig
	}

	return s.repository.FindPaginated(ctx, FindUserQuery{Country: req.Country}, req.Page, req.Size)
}

func (s *Service) ChangeEmail(ctx context.Context, userID uuid.UUID, new string) (User, error) {
	s.log.Infow("change email for user", "userID", userID)

	u, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		s.log.Errorw("failed to change email for user",
			"error", err,
			"userID", userID)
		return User{}, err
	}

	oldEmail := u.Email
	if u.Email == new {
		s.log.Infow("change email for user finished - emails are the same", "userID", userID)
		return u, nil
	}

	err = u.ChangeEmail(new)
	if err != nil {
		s.log.Errorw("failed to change email for user",
			"error", err,
			"userID", userID)
		return User{}, ErrActionFailed
	}

	err = s.repository.ChangeEmail(ctx, u)
	if err != nil {
		s.log.Errorw("failed to change email for user while persisting change in the repository",
			"error", err,
			"userID", userID)
		return User{}, err
	}

	evt := events.NewEvent(EvtUserEmailUpdated, EvtUserEmailUpdatedVersion, EvtUserEmailUpdatedPayload{
		ID:  userID,
		Old: oldEmail,
		New: new,
	})

	err = s.bus.Enqueue(evt)
	if err != nil {
		s.log.Errorw("failed to emit user email updated event",
			"error", err,
			"userID", userID)
		return User{}, ErrCreateFailed
	}

	return u, nil
}
