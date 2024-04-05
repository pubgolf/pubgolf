package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// VerifyPhoneNumber sets the player's phone number as verified, returning a boolean to indicate whether the phone number was previously unverified (i.e. false means the DB row was not updated).
func (q *Queries) VerifyPhoneNumber(ctx context.Context, num models.PhoneNum) (bool, error) {
	defer daoSpan(&ctx)()

	didUpdate, err := q.dbc.VerifyPhoneNumber(ctx, num)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("generate token: %w", err)
	}

	return didUpdate, nil
}
