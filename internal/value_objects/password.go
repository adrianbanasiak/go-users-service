package value_objects

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordMinLength = 15
)

var (
	ErrPasswordTooShort = errors.New("password is too short")
)

type Password struct {
	hash string
}

func (p *Password) String() string {
	return p.hash
}

func (p *Password) UnmarshalBSONValue(_ bsontype.Type, data []byte) error {
	s, _, ok := bsoncore.ReadString(data)
	if !ok {
		return errors.New("invalid data for Password")
	}

	err := ensureValidCountryCodeString(s)
	if err != nil {
		return err
	}

	p.hash = s

	return nil
}

func (p Password) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.String(p.hash).MarshalBSONValue()
}

func PasswordFromString(in string) (Password, error) {
	err := ensureValidPassword(in)
	if err != nil {
		return Password{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(password)}, nil
}

// ensureValidPassword ensures that given string matches standard for password security
func ensureValidPassword(in string) error {
	if len(in) < passwordMinLength {
		return ErrPasswordTooShort
	}

	// more validations

	return nil
}
