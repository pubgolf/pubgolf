package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// PlayerByPhoneNumber returns true if the phone number has been verified via an auth code.
func (q *Queries) PlayerByPhoneNumber(ctx context.Context, num models.PhoneNum) (models.Player, error) {
	defer daoSpan(&ctx)()

	pRow, err := q.dbc.PlayerByPhoneNumber(ctx, num)
	if err != nil {
		return models.Player{}, fmt.Errorf("fetch player: %w", err)
	}

	return models.Player{
		ID:   pRow.ID,
		Name: pRow.Name,
	}, nil
}
