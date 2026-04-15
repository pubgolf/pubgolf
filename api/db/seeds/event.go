package seeds

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// EventSummary holds the current state of an event's data for status reporting.
type EventSummary struct {
	Stages              int
	Rules               int
	Players             int
	Scores              int
	ScoresVerified      int
	ScoresUnverified    int
	AdjustmentTemplates int
	PlayerNames         []string
}

// PreviewEventData returns the current row counts and sample data for an event.
// The bool return indicates whether the event exists at all.
func PreviewEventData(ctx context.Context, db *sql.DB, eventKey string) (EventSummary, bool, error) {
	var eventID sql.NullString

	err := db.QueryRowContext(ctx,
		"SELECT id::text FROM events WHERE key = $1 AND deleted_at IS NULL", eventKey,
	).Scan(&eventID)
	if errors.Is(err, sql.ErrNoRows) || !eventID.Valid {
		return EventSummary{}, false, nil
	}

	if err != nil {
		return EventSummary{}, false, fmt.Errorf("check event existence: %w", err)
	}

	var s EventSummary

	queries := []struct {
		dest  *int
		query string
	}{
		{&s.Stages, "SELECT COUNT(*) FROM stages WHERE event_id = $1 AND deleted_at IS NULL"},
		{&s.Rules, `SELECT COUNT(DISTINCT r.id) FROM rules r
			JOIN stages s ON s.rule_id = r.id
			WHERE s.event_id = $1 AND s.deleted_at IS NULL AND r.deleted_at IS NULL`},
		{&s.Players, "SELECT COUNT(*) FROM event_players WHERE event_id = $1 AND deleted_at IS NULL"},
		{&s.Scores, `SELECT COUNT(*) FROM scores sc
			JOIN stages s ON sc.stage_id = s.id
			WHERE s.event_id = $1 AND sc.deleted_at IS NULL AND s.deleted_at IS NULL`},
		{&s.ScoresVerified, `SELECT COUNT(*) FROM scores sc
			JOIN stages s ON sc.stage_id = s.id
			WHERE s.event_id = $1 AND sc.is_verified = TRUE AND sc.deleted_at IS NULL AND s.deleted_at IS NULL`},
		{&s.AdjustmentTemplates, `SELECT COUNT(*) FROM adjustment_templates at
			LEFT JOIN stages s ON at.stage_id = s.id
			WHERE (at.event_id = $1 OR s.event_id = $1) AND at.deleted_at IS NULL`},
	}

	for _, q := range queries {
		scanErr := db.QueryRowContext(ctx, q.query, eventID.String).Scan(q.dest)
		if scanErr != nil {
			return EventSummary{}, true, fmt.Errorf("query event data: %w", scanErr)
		}
	}

	s.ScoresUnverified = s.Scores - s.ScoresVerified

	// Sample player names.
	rows, err := db.QueryContext(ctx, `
		SELECT p.name FROM players p
		JOIN event_players ep ON ep.player_id = p.id
		WHERE ep.event_id = $1 AND ep.deleted_at IS NULL AND p.deleted_at IS NULL
		ORDER BY p.name
		LIMIT 5`, eventID.String)
	if err != nil {
		return EventSummary{}, true, fmt.Errorf("query player names: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		scanErr := rows.Scan(&name)
		if scanErr != nil {
			return EventSummary{}, true, fmt.Errorf("scan player name: %w", scanErr)
		}

		s.PlayerNames = append(s.PlayerNames, name)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return EventSummary{}, true, fmt.Errorf("iterate player names: %w", rowsErr)
	}

	return s, true, nil
}

// DeleteEventData removes all data associated with an event key in FK-safe order.
// Players and venues are left intact as shared resources.
func DeleteEventData(ctx context.Context, tx *sql.Tx, eventKey string) error {
	// Resolve event ID first.
	var eventID string

	err := tx.QueryRowContext(ctx,
		"SELECT id::text FROM events WHERE key = $1", eventKey,
	).Scan(&eventID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil // Nothing to delete.
	}

	if err != nil {
		return fmt.Errorf("resolve event ID: %w", err)
	}

	// Delete in FK order. All queries take eventID as $1.
	deletes := []string{
		`DELETE FROM adjustments WHERE stage_id IN (
			SELECT id FROM stages WHERE event_id = $1)`,
		`DELETE FROM scores WHERE stage_id IN (
			SELECT id FROM stages WHERE event_id = $1)`,
		`DELETE FROM auth_tokens WHERE player_id IN (
			SELECT player_id FROM event_players WHERE event_id = $1)`,
		"DELETE FROM event_players WHERE event_id = $1",
		"DELETE FROM adjustment_templates WHERE event_id = $1",
		`DELETE FROM adjustment_templates WHERE stage_id IN (
			SELECT id FROM stages WHERE event_id = $1)`,
		"DELETE FROM stages WHERE event_id = $1",
		"DELETE FROM events WHERE id = $1",
	}

	for _, q := range deletes {
		_, execErr := tx.ExecContext(ctx, q, eventID)
		if execErr != nil {
			return fmt.Errorf("delete event data: %w", execErr)
		}
	}

	// Delete orphaned rules (rules no longer referenced by any stage).
	_, execErr := tx.ExecContext(ctx,
		`DELETE FROM rules WHERE deleted_at IS NULL
			AND id NOT IN (SELECT DISTINCT rule_id FROM stages WHERE rule_id IS NOT NULL)`)
	if execErr != nil {
		return fmt.Errorf("delete orphaned rules: %w", execErr)
	}

	return nil
}
