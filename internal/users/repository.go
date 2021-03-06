package users

import (
	"context"
	"github.com/google/uuid"
)

type FindUserQuery struct {
	Country string
}

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	ChangeEmail(ctx context.Context, user User) error
	Delete(ctx context.Context, ID uuid.UUID) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindPaginated(ctx context.Context, query FindUserQuery, page, limit int) ([]User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (User, error)
}
