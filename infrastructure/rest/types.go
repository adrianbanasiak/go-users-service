package rest

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	Successful bool
	Data       any
	Errors     []string
	SentAt     time.Time
}

type CreateUserReq struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nickname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
}

type UserRes struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	NickName    string    `json:"nickname"`
	Email       string    `json:"email"`
	CountryCode string    `json:"country"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
