package users

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, ID uuid.UUID) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindPaginated(ctx context.Context, page, limit int) ([]User, error)
}
