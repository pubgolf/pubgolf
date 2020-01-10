package utils

import (
	"regexp"
)

var uuidFormat *regexp.Regexp = regexp.MustCompile(
	"^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
var phoneNumberFormat *regexp.Regexp = regexp.MustCompile(
	"^\\+[1-9]\\d{1,14}$")

// IsEmpty returns true if a string is empty.
func IsEmpty(arg *string) bool {
	return *arg == ""
}

// InvalidPhoneNumberFormat returns true if a phone number string doesn't match the E.164 format:
// https://www.twilio.com/docs/glossary/what-e164
func InvalidPhoneNumberFormat(phoneNumber *string) bool {
	isValid := phoneNumberFormat.MatchString(*phoneNumber)
	return !isValid
}

// InvalidIDFormat returns true if the provided ID is not a valid UUID v4.
func InvalidIDFormat(id *string) bool {
	isValid := uuidFormat.MatchString(*id)
	return !isValid
}

// InvalidAuthCodeFormat returns true if the auth code is not a number in the valid range.
func InvalidAuthCodeFormat(authCode uint32) bool {
	return authCode < 100000 || authCode > 999999
}
