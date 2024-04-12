package dao

import (
	"context"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// UpsertRegistration creates a new player and adds them to the given event.
func (q *Queries) UpsertRegistration(ctx context.Context, playerID models.PlayerID, eventKey string, cat models.ScoringCategory) error {
	defer daoSpan(&ctx)()

	eventID, err := q.dbc.EventIDByKey(ctx, eventKey)
	if err != nil {
		return fmt.Errorf("lookup event ID: %w", err)
	}

	err = q.dbc.UpsertRegistration(ctx, dbc.UpsertRegistrationParams{
		PlayerID:        playerID,
		EventID:         eventID,
		ScoringCategory: cat,
	})
	if err != nil {
		return fmt.Errorf("upsert registration: %w", err)
	}

	return nil
}
