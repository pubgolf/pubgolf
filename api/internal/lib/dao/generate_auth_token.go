package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// GenerateAuthToken generates an auth token for the player with the given phone number, returning the auth token and the player's ID.
func (q *Queries) GenerateAuthToken(ctx context.Context, num models.PhoneNum) (models.PlayerID, models.AuthToken, error) {
	defer daoSpan(&ctx)()

	res, err := q.dbc.GenerateAuthToken(ctx, num)
	if err != nil {
		return models.PlayerID{}, models.AuthToken{}, fmt.Errorf("generate token: %w", err)
	}

	return res.PlayerID, res.AuthToken, nil
}
