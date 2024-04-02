package models

import (
	"errors"
	"fmt"
	"regexp"
)

// PhoneNumPattern validates that a string matches the E.164 format: https://www.twilio.com/docs/glossary/what-e164
var PhoneNumPattern = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// ErrInvalidPhoneNumFormat is returned if a phone number can't be parsed due to invalid formatting.
var ErrInvalidPhoneNumFormat = errors.New("not a valid phone number in E.164 format")

type PhoneNum string

// NewPhoneNum attempts to parse a phone number as E.164 format.
func NewPhoneNum(num string) (PhoneNum, error) {
	if !PhoneNumPattern.MatchString(num) {
		return "", fmt.Errorf("parse phone number %q: %w", num, ErrInvalidPhoneNumFormat)
	}

	return PhoneNum(num), nil
}

// String returns a string representation of the phone number in E.164 format.
func (pn PhoneNum) String() string {
	return string(pn)
}
