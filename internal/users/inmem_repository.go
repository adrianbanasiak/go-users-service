package users

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

func NewInMemoryRepository(wantErr error) *InMemoryRepository {
	return &InMemoryRepository{wantErr: wantErr}
}

type InMemoryRepository struct {
	users   []User
	wantErr error
	called  bool
}

func (r *InMemoryRepository) ChangeEmail(_ context.Context, _ User) error {
	panic("implement me")
}

func (r *InMemoryRepository) FindByID(_ context.Context, userID uuid.UUID) (User, error) {
	if r.wantErr != nil {
		return User{}, r.wantErr
	}

	for _, user := range r.users {
		if user.ID == userID {
			return user, nil
		}
	}

	return User{}, ErrNotFound
}

func (r *InMemoryRepository) FindPaginated(_ context.Context, _, _ int) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *InMemoryRepository) FindByEmail(_ context.Context, email string) (User, error) {
	if r.wantErr != nil {
		return User{}, r.wantErr
	}

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotFound
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
