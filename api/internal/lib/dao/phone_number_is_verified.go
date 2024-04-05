package dao

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PhoneNumberIsVerified returns true if the phone number has been verified via an auth code.
func (q *Queries) PhoneNumberIsVerified(ctx context.Context, num models.PhoneNum) (bool, error) {
	defer daoSpan(&ctx)()

	return q.dbc.PhoneNumberIsVerified(ctx, num) //nolint:wrapcheck // Trivial passthrough
}
