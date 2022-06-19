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
		ID:          uuid.New(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		NickName:    req.NickName,
		password:    p,
		Email:       req.Email,
		CountryCode: cc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	NickName    string
	password    value_objects.Password
	Email       string
	CountryCode value_objects.CountryCode
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
