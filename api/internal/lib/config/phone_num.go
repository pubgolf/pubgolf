package config

import (
	"fmt"
	"strings"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PhoneNumSet contains a list of phone numbers to match against, or, if universal is set to true, will match any valid phone number.
type PhoneNumSet struct {
	universal bool
	values    map[models.PhoneNum]struct{}
}

// Set parses a PhoneNumSet from a comma-separated list of E.164 formatted phone numbers, or treats the literal "*" as a wildcard. Extra commas will be ignored, and an error will be thrown on any malformed phone numbers.
func (pns *PhoneNumSet) Set(value string) error {
	if value == "*" {
		*pns = PhoneNumSet{
			universal: true,
			values:    nil,
		}

		return nil
	}

	ns := make(map[models.PhoneNum]struct{})

	for _, n := range strings.Split(value, ",") {
		n = strings.TrimSpace(n)

		if n == "" {
			continue
		}

		num, err := models.NewPhoneNum(n)
		if err != nil {
			return fmt.Errorf("create phone num set: %w", err)
		}

		ns[num] = struct{}{}
	}

	*pns = PhoneNumSet{
		universal: false,
		values:    ns,
	}

	return nil
}

// Match returns true if the provided num is contained within the set.
func (pns PhoneNumSet) Match(num models.PhoneNum) bool {
	if pns.universal {
		return true
	}

	_, ok := pns.values[num]

	return ok
}
