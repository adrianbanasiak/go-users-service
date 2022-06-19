package users

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, ID uuid.UUID) error
	List(ctx context.Context, page, items int) ([]User, error)
}
