package value_objects

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var (
	ccRegxp               = regexp.MustCompile("^[a-zA-Z]{2}$")
	ErrInvalidCountryCode = errors.New("invalid country code value provided")
)

type CountryCode struct {
	value string
}

func (code *CountryCode) String() string {
	return code.value
}

func (code *CountryCode) UnmarshalJSON(bytes []byte) error {
	var s string
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	err = ensureValidCountryCodeString(s)
	if err != nil {
		return err
	}

	code.value = strings.ToUpper(s)

	return nil
}

func (code CountryCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(code.value)
}

func CountryCodeFromString(in string) (CountryCode, error) {
	if err := ensureValidCountryCodeString(in); err != nil {
		return CountryCode{}, ErrInvalidCountryCode
	}

	return CountryCode{value: strings.ToUpper(in)}, nil
}

func ensureValidCountryCodeString(in string) error {
	if !ccRegxp.MatchString(in) {
		return ErrInvalidCountryCode
	}
	return nil
}
