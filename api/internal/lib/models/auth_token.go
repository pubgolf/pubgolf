package models

// RedactedString returns a redacted representation usable for logging.
func (at AuthToken) RedactedString() string {
	s := at.String()

	return s[:3] + "********************" + s[len(s)-3:]
}
