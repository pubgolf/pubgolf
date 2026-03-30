package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

// testClients holds the RPC client pair used by every E2E test.
type testClients struct {
	pub   apiv1connect.PubGolfServiceClient
	admin apiv1connect.AdminServiceClient
}

func newTestClients() testClients {
	return testClients{
		pub:   apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc"),
		admin: apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc"),
	}
}

// seededEvent holds IDs returned after seeding an event with venues and stages.
type seededEvent struct {
	eventID  models.EventID
	stageIDs []models.StageID // len == numStages, ordered by rank
}

// seedEvent inserts an event with the given key and startsAtExpr (a SQL expression
// relative to NOW(), e.g. "NOW() + '30 minutes'" or "NOW() + '-45 minutes'"),
// then inserts numStages venues+stages with 30-minute durations and ranks 10,20,...
//
// It purges all caches after seeding.
func seedEvent(ctx context.Context, t *testing.T, db *sql.DB, tc testClients, eventKey string, startsAtExpr string, numStages int) seededEvent {
	t.Helper()

	row := db.QueryRowContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, "+startsAtExpr+") RETURNING id", eventKey) //nolint:gosec // startsAtExpr is always a SQL literal from test code, not user input
	require.NoError(t, row.Err(), "seed DB: insert event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	stageIDs := make([]models.StageID, numStages)

	for i := range numStages {
		row := db.QueryRowContext(ctx, "INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id",
			fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), "seed DB: insert venue %d", i)

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		row = db.QueryRowContext(ctx,
			"INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30) RETURNING id",
			eventID, venueID, (i+1)*10)
		require.NoError(t, row.Err(), "seed DB: insert stage %d", i)
		require.NoError(t, row.Scan(&stageIDs[i]), "scan returned stage ID %d", i)
	}

	_, err := tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err, "purge caches after seeding event")

	return seededEvent{eventID: eventID, stageIDs: stageIDs}
}

// seededPlayer holds IDs returned after creating a player with an auth token.
type seededPlayer struct {
	id    models.PlayerID
	token string
}

// seedPlayer creates a player via the admin API and inserts an auth token row.
// Returns the player ID and auth token string ready for use with requestWithAuth.
// Note: does not purge caches. Callers that did not use seedEvent should call
// PurgeAllCaches before making public API requests.
func seedPlayer(ctx context.Context, t *testing.T, db *sql.DB, tc testClients, phone string, eventKey string, category apiv1.ScoringCategory, name string) seededPlayer {
	t.Helper()

	playerResp, err := tc.admin.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: name,
		},
		PhoneNumber: phone,
		Registration: &apiv1.EventRegistration{
			EventKey:        eventKey,
			ScoringCategory: category,
		},
	}))
	require.NoError(t, err, "create player %s", phone)

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID %s", phone)

	row := db.QueryRowContext(ctx, "INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token for %s", phone)

	var tok models.AuthToken
	require.NoError(t, row.Scan(&tok), "scan auth token for %s", phone)

	return seededPlayer{id: playerID, token: tok.String()}
}
