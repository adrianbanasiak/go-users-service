package value_objects

import (
	"encoding/json"
	"errors"
)

const (
	passwordMinLength = 15
	marshalString     = "****PASSWORD*****"
)

var (
	ErrPasswordTooShort = errors.New("password is too short")
)

type Password struct {
	value string
}

func (p *Password) UnmarshalJSON(bytes []byte) error {
	var s string
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	err = ensureValidPassword(s)
	if err != nil {
		return err
	}

	p.value = s
	return nil
}

func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshalString)
}

func PasswordFromString(in string) (Password, error) {
	err := ensureValidPassword(in)
	if err != nil {
		return Password{}, err
	}

	return Password{value: in}, nil
}

// ensureValidPassword ensures that given string matches standard for password security
func ensureValidPassword(in string) error {
	if len(in) < passwordMinLength {
		return ErrPasswordTooShort
	}

	// more validations

	return nil
}
