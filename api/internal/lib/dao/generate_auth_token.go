package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// GenerateAuthTokenResult contains the generated auth token and affected player from a GenerateAuthToken call.
type GenerateAuthTokenResult struct {
	// DidInvalidate indicates whether there were existing auth tokens which were deactivated.
	DidInvalidate bool
	PlayerID      models.PlayerID
	AuthToken     models.AuthToken
}

// GenerateAuthToken generates an auth token for the player with the given phone number, returning the auth token and the player's ID.
func (q *Queries) GenerateAuthToken(ctx context.Context, num models.PhoneNum) (GenerateAuthTokenResult, error) {
	defer daoSpan(&ctx)()

	gat := GenerateAuthTokenResult{}
	err := q.useTx(ctx, func(ctx context.Context, q *Queries) error {
		didUpdate, err := q.dbc.DeactivateAuthTokens(ctx, num)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("deactivate old token: %w", err)
		}

		gat.DidInvalidate = didUpdate

		res, err := q.dbc.GenerateAuthToken(ctx, num)
		if err != nil {
			return fmt.Errorf("generate new token: %w", err)
		}

		gat.PlayerID = res.PlayerID
		gat.AuthToken = res.AuthToken

		return nil
	})

	return gat, err
}
