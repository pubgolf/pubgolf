package seeds

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const adminTestEventKey = "nyc-test"

// AdminTestSeed populates a local or remote database with a realistic test event.
var AdminTestSeed = Seed{
	Name:     "admin-test",
	EventKey: adminTestEventKey,
	Expected: ExpectedCounts{
		Stages:              9,
		Rules:               9,
		Players:             15,
		Scores:              30,
		AdjustmentTemplates: 5,
	},
	Run: runAdminTest,
}

type venue struct {
	name    string
	address string
}

type player struct {
	name  string
	phone string
	// scoringCategory is the enum value from enum_scoring_categories.
	scoringCategory string
}

var adminTestVenues = []venue{
	{"The Dead Rabbit", "30 Water St, New York, NY 10004"},
	{"McSorley's Old Ale House", "15 E 7th St, New York, NY 10003"},
	{"Blind Tiger Ale House", "281 Bleecker St, New York, NY 10014"},
	{"Ear Inn", "326 Spring St, New York, NY 10013"},
	{"Pete's Tavern", "129 E 18th St, New York, NY 10003"},
	{"Old Town Bar", "45 E 18th St, New York, NY 10003"},
	{"Muldoon's", "692 3rd Ave, New York, NY 10017"},
	{"The Long Hall", "54 W 36th St, New York, NY 10018"},
	{"Paddy Reilly's Music Bar", "519 2nd Ave, New York, NY 10016"},
}

var adminTestRuleDescriptions = []string{
	"Order a craft beer on tap. Bonus for trying something new.",
	"Take a photo with the bartender. Must be posted to the group chat.",
	"Sing along to any song playing for at least 30 seconds.",
	"Play a round of darts or pool if available. Otherwise, thumb wrestle.",
	"Order a drink you've never had before and rate it 1-10.",
	"Find and photograph the oldest item on the wall.",
	"Make a toast to the group. Must be at least 3 sentences.",
	"Order the cheapest drink on the menu. Document the price.",
	"Challenge someone outside your group to a game or conversation.",
}

var adminTestPlayers = []player{
	// 9-hole (5 players)
	{"Alice Johnson", "+15550001001", "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE"},
	{"Bob Smith", "+15550001002", "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE"},
	{"Charlie Brown", "+15550001003", "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE"},
	{"Diana Prince", "+15550001004", "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE"},
	{"Eve Martinez", "+15550001005", "SCORING_CATEGORY_PUB_GOLF_NINE_HOLE"},
	// 5-hole (5 players)
	{"Frank Ocean", "+15550001006", "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE"},
	{"Grace Lee", "+15550001007", "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE"},
	{"Hank Williams", "+15550001008", "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE"},
	{"Iris Chang", "+15550001009", "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE"},
	{"Jack Kerouac", "+15550001010", "SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE"},
	// Challenges (5 players)
	{"Karen Wu", "+15550001011", "SCORING_CATEGORY_PUB_GOLF_CHALLENGES"},
	{"Leo Tolstoy", "+15550001012", "SCORING_CATEGORY_PUB_GOLF_CHALLENGES"},
	{"Mia Farrow", "+15550001013", "SCORING_CATEGORY_PUB_GOLF_CHALLENGES"},
	{"Nick Cave", "+15550001014", "SCORING_CATEGORY_PUB_GOLF_CHALLENGES"},
	{"Olivia Wilde", "+15550001015", "SCORING_CATEGORY_PUB_GOLF_CHALLENGES"},
}

func runAdminTest(ctx context.Context, tx *sql.Tx) error {
	// 1. Insert event.
	var eventID string

	err := tx.QueryRowContext(ctx,
		"INSERT INTO events (key, starts_at) VALUES ($1, NOW() + INTERVAL '1 day') RETURNING id::text",
		adminTestEventKey,
	).Scan(&eventID)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	// 2. Insert venues (upsert — venues are shared resources).
	venueIDs := make([]string, len(adminTestVenues))
	for i, v := range adminTestVenues {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO venues (name, address)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			RETURNING id::text`,
			v.name, v.address,
		).Scan(&venueIDs[i])
		if errors.Is(err, sql.ErrNoRows) {
			// Already exists, look it up.
			err = tx.QueryRowContext(ctx,
				"SELECT id::text FROM venues WHERE name = $1 AND address = $2",
				v.name, v.address,
			).Scan(&venueIDs[i])
		}

		if err != nil {
			return fmt.Errorf("insert venue %q: %w", v.name, err)
		}
	}

	// 3. Insert rules.
	ruleIDs := make([]string, len(adminTestRuleDescriptions))
	for i, r := range adminTestRuleDescriptions {
		err := tx.QueryRowContext(ctx,
			"INSERT INTO rules (description) VALUES ($1) RETURNING id::text",
			r,
		).Scan(&ruleIDs[i])
		if err != nil {
			return fmt.Errorf("insert rule %d: %w", i+1, err)
		}
	}

	// 4. Insert stages (one per venue+rule pair, ranked 10..90).
	stageIDs := make([]string, len(adminTestVenues))
	for i := range adminTestVenues {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO stages (event_id, venue_id, rule_id, rank, duration_minutes)
			VALUES ($1, $2, $3, $4, 30)
			RETURNING id::text`,
			eventID, venueIDs[i], ruleIDs[i], (i+1)*10,
		).Scan(&stageIDs[i])
		if err != nil {
			return fmt.Errorf("insert stage %d: %w", i+1, err)
		}
	}

	// 5. Insert players (upsert on phone_number — players are shared resources).
	playerIDs := make([]string, len(adminTestPlayers))
	for i, p := range adminTestPlayers {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO players (name, phone_number, phone_number_verified)
			VALUES ($1, $2, TRUE)
			ON CONFLICT (phone_number) DO UPDATE SET name = EXCLUDED.name
			RETURNING id::text`,
			p.name, p.phone,
		).Scan(&playerIDs[i])
		if err != nil {
			return fmt.Errorf("insert player %q: %w", p.name, err)
		}
	}

	// 6. Insert event_players.
	for i, p := range adminTestPlayers {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO event_players (event_id, player_id, scoring_category)
			VALUES ($1, $2, $3)`,
			eventID, playerIDs[i], p.scoringCategory,
		)
		if err != nil {
			return fmt.Errorf("insert event_player %q: %w", p.name, err)
		}
	}

	// 7. Insert auth_tokens (one per player).
	for i, p := range adminTestPlayers {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO auth_tokens (player_id) VALUES ($1)",
			playerIDs[i],
		)
		if err != nil {
			return fmt.Errorf("insert auth_token for %q: %w", p.name, err)
		}
	}

	// 8. Insert scores — 2 per player (stages 0 and 1), alternating verified.
	for i := range adminTestPlayers {
		for j := range 2 {
			score := 3 + (i+j)%5 // Values 3-7 for variety.
			isVerified := (i+j)%3 != 0

			_, err := tx.ExecContext(ctx, `
				INSERT INTO scores (stage_id, player_id, value, is_verified)
				VALUES ($1, $2, $3, $4)`,
				stageIDs[j], playerIDs[i], score, isVerified,
			)
			if err != nil {
				return fmt.Errorf("insert score (player %d, stage %d): %w", i+1, j+1, err)
			}
		}
	}

	// 9. Insert adjustment_templates.
	// 3 global (event-scoped), 2 stage-specific.
	// CHECK constraint: exactly one of event_id or stage_id must be set.
	globalTemplates := []struct {
		label string
		value int
		rank  int
	}{
		{"Late arrival", -1, 1},         // visible
		{"Best dressed", 2, 2},          // visible
		{"Organizer discretion", -3, 0}, // hidden (rank=0)
	}

	for _, t := range globalTemplates {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO adjustment_templates (event_id, value, label, rank)
			VALUES ($1, $2, $3, $4)`,
			eventID, t.value, t.label, t.rank,
		)
		if err != nil {
			return fmt.Errorf("insert global adjustment template %q: %w", t.label, err)
		}
	}

	stageTemplates := []struct {
		stageIdx int
		label    string
		value    int
		rank     int
	}{
		{0, "Photo bonus", 1, 1},        // visible, stage 1
		{2, "Challenge penalty", -2, 0}, // hidden, stage 3
	}

	for _, t := range stageTemplates {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO adjustment_templates (stage_id, value, label, rank)
			VALUES ($1, $2, $3, $4)`,
			stageIDs[t.stageIdx], t.value, t.label, t.rank,
		)
		if err != nil {
			return fmt.Errorf("insert stage adjustment template %q: %w", t.label, err)
		}
	}

	return nil
}
