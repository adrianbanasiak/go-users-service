package rest

import (
	"github.com/adrianbanasiak/go-users-service/internal/users"
)

func PresentUser(user users.User) UserRes {
	return UserRes{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		NickName:    user.NickName,
		Email:       user.Email,
		CountryCode: user.CountryCode.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
