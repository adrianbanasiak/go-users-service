package users

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	users   []User
	wantErr error
	called  bool
}

func NewInMemoryRepository(wantErr error) *InMemoryRepository {
	return &InMemoryRepository{wantErr: wantErr}
}

func (r *InMemoryRepository) total() int {
	return len(r.users)
}

func (r *InMemoryRepository) Called() bool {
	return r.called
}

func (r *InMemoryRepository) Create(_ context.Context, user User) (User, error) {
	r.called = true

	if r.wantErr != nil {
		return User{}, r.wantErr
	}

	r.users = append(r.users, user)

	return user, nil
}

func (r *InMemoryRepository) Delete(_ context.Context, ID uuid.UUID) error {
	r.called = true

	if r.wantErr != nil {
		return r.wantErr
	}

	for i, user := range r.users {
		if user.ID != ID {
			continue
		}

		r.users = append(r.users[:i], r.users[i+1:]...)
		return nil
	}

	return errors.New("not found")
}

func (r *InMemoryRepository) List(_ context.Context, _, _ int) ([]User, error) {
	return r.users, nil
}
