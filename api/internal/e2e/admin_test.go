//nolint:paralleltest
package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1/apiv1connect"
)

func Test_AdminAdjustmentTemplates(t *testing.T) {
	testEventKey := "test-event-key-admin-adj"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event + venues/stages.
	row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '-45 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	for i := range 9 {
		row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), "seed DB: insert venue %d", i)

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30)", eventID, venueID, (i+1)*10)
		require.NoError(t, err, "seed DB: insert stage %d", i)
	}

	_, err := ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// Create an adjustment template.
	_, err = ac.CreateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: testEventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Test Penalty",
				Value: 2,
			},
			Rank:      10,
			IsVisible: true,
		},
	}))
	require.NoError(t, err, "CreateAdjustmentTemplate")

	// List adjustment templates and verify the created one is present.
	list, err := ac.ListAdjustmentTemplates(ctx, requestWithAdminAuth(&apiv1.ListAdjustmentTemplatesRequest{
		EventKey: testEventKey,
	}))
	require.NoError(t, err, "ListAdjustmentTemplates")

	templates := list.Msg.GetTemplates()
	require.NotEmpty(t, templates, "at least one template exists")

	var foundTemplate *apiv1.AdjustmentTemplate

	for _, tmpl := range templates {
		if tmpl.GetData().GetAdjustment().GetLabel() == "Test Penalty" {
			foundTemplate = tmpl

			break
		}
	}

	require.NotNil(t, foundTemplate, "created template found in list")
	assert.EqualValues(t, 2, foundTemplate.GetData().GetAdjustment().GetValue(), "template value matches")
	assert.EqualValues(t, 10, foundTemplate.GetData().GetRank(), "template rank matches")

	templateID := foundTemplate.GetId()

	// Update the adjustment template.
	_, err = ac.UpdateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.UpdateAdjustmentTemplateRequest{
		Template: &apiv1.AdjustmentTemplate{
			Id: templateID,
			Data: &apiv1.AdjustmentTemplateData{
				EventKey: testEventKey,
				Adjustment: &apiv1.AdjustmentData{
					Label: "Updated Penalty",
					Value: 5,
				},
				Rank:      20,
				IsVisible: true,
			},
		},
	}))
	require.NoError(t, err, "UpdateAdjustmentTemplate")

	// Verify update is reflected.
	list, err = ac.ListAdjustmentTemplates(ctx, requestWithAdminAuth(&apiv1.ListAdjustmentTemplatesRequest{
		EventKey: testEventKey,
	}))
	require.NoError(t, err, "ListAdjustmentTemplates after update")

	var updatedTemplate *apiv1.AdjustmentTemplate

	for _, tmpl := range list.Msg.GetTemplates() {
		if tmpl.GetId() == templateID {
			updatedTemplate = tmpl

			break
		}
	}

	require.NotNil(t, updatedTemplate, "updated template found in list")
	assert.Equal(t, "Updated Penalty", updatedTemplate.GetData().GetAdjustment().GetLabel(), "label updated")
	assert.EqualValues(t, 5, updatedTemplate.GetData().GetAdjustment().GetValue(), "value updated")
	assert.EqualValues(t, 20, updatedTemplate.GetData().GetRank(), "rank updated")
}

func Test_AdminScoreManagement(t *testing.T) {
	testEventKey := "test-event-key-admin-score"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")
	c := apiv1connect.NewPubGolfServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event started 45min ago + venues/stages.
	row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '-45 minutes') RETURNING id", testEventKey)
	require.NoError(t, row.Err(), "seed DB: insert event")

	var eventID models.EventID
	require.NoError(t, row.Scan(&eventID), "scan returned event ID")

	var secondStageID models.StageID

	for i := range 9 {
		row := sharedTestDB.QueryRowContext(ctx, "INSERT INTO venues (name, address) VALUES ($1, $2) RETURNING id", fmt.Sprintf("Venue %d", i+1), fmt.Sprintf("%d Test St", i+1))
		require.NoError(t, row.Err(), "seed DB: insert venue %d", i)

		var venueID models.VenueID
		require.NoError(t, row.Scan(&venueID), "scan returned venue ID")

		row = sharedTestDB.QueryRowContext(ctx, "INSERT INTO stages (event_id, venue_id, rank, duration_minutes) VALUES ($1, $2, $3, 30) RETURNING id", eventID, venueID, (i+1)*10)
		require.NoError(t, row.Err(), "seed DB: insert stage %d", i)

		if i == 1 {
			require.NoError(t, row.Scan(&secondStageID), "scan returned stage ID")
		}
	}

	// Create player.
	playerResp, err := ac.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "ScoreTestPlayer",
		},
		PhoneNumber: "+15559380301",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		},
	}))
	require.NoError(t, err, "create player")

	playerID, err := models.PlayerIDFromString(playerResp.Msg.GetPlayer().GetId())
	require.NoError(t, err, "convert player ID")

	row = sharedTestDB.QueryRowContext(ctx, "INSERT INTO auth_tokens (player_id) VALUES ($1) RETURNING id", playerID)
	require.NoError(t, row.Err(), "insert auth token")

	var playerToken models.AuthToken
	require.NoError(t, row.Scan(&playerToken), "scan returned auth token")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// CreateStageScore via admin.
	createRes, err := ac.CreateStageScore(ctx, requestWithAdminAuth(&apiv1.CreateStageScoreRequest{
		Data: &apiv1.StageScoreData{
			StageId:  secondStageID.String(),
			PlayerId: playerID.String(),
			Score: &apiv1.ScoreData{
				Value: 5,
			},
		},
	}))
	require.NoError(t, err, "CreateStageScore")
	assert.NotNil(t, createRes.Msg.GetScore(), "created score is returned")

	// Verify via GetScoresForPlayer.
	scores, err := c.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: testEventKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err, "GetScoresForPlayer after create")

	scoreEntries := scores.Msg.GetScoreBoard().GetScores()
	require.NotEmpty(t, scoreEntries, "has score entries after create")

	// UpdateStageScore via admin.
	updateRes, err := ac.UpdateStageScore(ctx, requestWithAdminAuth(&apiv1.UpdateStageScoreRequest{
		Score: &apiv1.StageScore{
			StageId:  secondStageID.String(),
			PlayerId: playerID.String(),
			Score: &apiv1.Score{
				Data: &apiv1.ScoreData{
					Value: 7,
				},
			},
			IsVerified: true,
		},
	}))
	require.NoError(t, err, "UpdateStageScore")
	assert.NotNil(t, updateRes.Msg.GetScore(), "updated score is returned")

	// Verify updated score via GetScoresForPlayer.
	scores, err = c.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: testEventKey,
		PlayerId: playerID.String(),
	}, playerToken.String()))
	require.NoError(t, err, "GetScoresForPlayer after update")
	require.NotEmpty(t, scores.Msg.GetScoreBoard().GetScores(), "has score entries after update")

	// DeleteStageScore via admin.
	_, err = ac.DeleteStageScore(ctx, requestWithAdminAuth(&apiv1.DeleteStageScoreRequest{
		StageId:  secondStageID.String(),
		PlayerId: playerID.String(),
	}))
	require.NoError(t, err, "DeleteStageScore")

	// Verify removal via ListStageScores.
	listRes, err := ac.ListStageScores(ctx, requestWithAdminAuth(&apiv1.ListStageScoresRequest{
		EventKey: testEventKey,
	}))
	require.NoError(t, err, "ListStageScores after delete")

	for _, s := range listRes.Msg.GetScores() {
		assert.False(t,
			s.GetStageId() == secondStageID.String() && s.GetPlayerId() == playerID.String(),
			"deleted score should not appear in ListStageScores",
		)
	}
}

func Test_AdminPlayerManagement(t *testing.T) {
	testEventKey := "test-event-key-admin-player"
	ctx := t.Context()

	ac := apiv1connect.NewAdminServiceClient(http.DefaultClient, "http://localhost:3100/rpc")

	// Insert event.
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", testEventKey)
	require.NoError(t, err, "seed DB: insert future event")

	_, err = ac.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// CreatePlayer via admin with registration.
	playerResp, err := ac.CreatePlayer(ctx, requestWithAdminAuth(&apiv1.AdminServiceCreatePlayerRequest{
		PlayerData: &apiv1.PlayerData{
			Name: "AdminCreated",
		},
		PhoneNumber: "+15559380302",
		Registration: &apiv1.EventRegistration{
			EventKey:        testEventKey,
			ScoringCategory: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		},
	}))
	require.NoError(t, err, "CreatePlayer")

	createdPlayerID := playerResp.Msg.GetPlayer().GetId()
	require.NotEmpty(t, createdPlayerID, "created player has ID")

	// ListPlayers and verify the created player is present.
	list, err := ac.ListPlayers(ctx, requestWithAdminAuth(&apiv1.ListPlayersRequest{
		EventKey: testEventKey,
	}))
	require.NoError(t, err, "ListPlayers")

	var found bool

	for _, p := range list.Msg.GetPlayers() {
		if p.GetId() == createdPlayerID {
			found = true

			assert.Equal(t, "AdminCreated", p.GetData().GetName(), "player name matches")

			break
		}
	}

	assert.True(t, found, "created player found in ListPlayers")
}
