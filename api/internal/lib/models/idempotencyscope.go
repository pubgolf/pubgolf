package models

// IdempotencyScope describes the operation type for an idempotency key.
//
//nolint:recvcheck // TODO: Remove this once we figure out how to properly exclude generated files
type IdempotencyScope int

// IdempotencyScope values.
const (
	IdempotencyScopeScoreSubmission IdempotencyScope = iota
)
