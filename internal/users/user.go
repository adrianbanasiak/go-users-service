package users

import (
	"github.com/adrianbanasiak/go-users-service/internal/value_objects"
	"github.com/google/uuid"
	"time"
)

func NewUser(req CreateUserReq) (User, error) {
	p, err := value_objects.PasswordFromString(req.Password)
	if err != nil {
		return User{}, err
	}

	cc, err := value_objects.CountryCodeFromString(req.CountryCode)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:           uuid.New(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		NickName:     req.NickName,
		PasswordHash: p.String(),
		Email:        req.Email,
		CountryCode:  cc,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

type User struct {
	ID           uuid.UUID                 `bson:"_id"`
	FirstName    string                    `bson:"first_name"`
	LastName     string                    `bson:"last_name"`
	NickName     string                    `bson:"nick_name"`
	PasswordHash string                    `bson:"password_hash"`
	Email        string                    `bson:"email"`
	CountryCode  value_objects.CountryCode `bson:"country_code"`
	CreatedAt    time.Time                 `bson:"created_at"`
	UpdatedAt    time.Time                 `bson:"updated_at"`
}
