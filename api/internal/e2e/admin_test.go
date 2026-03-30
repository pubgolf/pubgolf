//nolint:paralleltest
package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_AdminAdjustmentTemplates(t *testing.T) {
	const eventKey = "test-event-key-admin-adj"

	ctx := t.Context()
	tc := newTestClients()

	seedEvent(ctx, t, sharedTestDB, tc, eventKey, "NOW() + '-45 minutes'", 9)

	// Create an adjustment template.
	_, err := tc.admin.CreateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.CreateAdjustmentTemplateRequest{
		Data: &apiv1.AdjustmentTemplateData{
			EventKey: eventKey,
			Adjustment: &apiv1.AdjustmentData{
				Label: "Test Penalty",
				Value: 2,
			},
			Rank:      10,
			IsVisible: true,
		},
	}))
	require.NoError(t, err, "CreateAdjustmentTemplate")

	// List and verify the created template.
	list, err := tc.admin.ListAdjustmentTemplates(ctx, requestWithAdminAuth(&apiv1.ListAdjustmentTemplatesRequest{
		EventKey: eventKey,
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
	_, err = tc.admin.UpdateAdjustmentTemplate(ctx, requestWithAdminAuth(&apiv1.UpdateAdjustmentTemplateRequest{
		Template: &apiv1.AdjustmentTemplate{
			Id: templateID,
			Data: &apiv1.AdjustmentTemplateData{
				EventKey: eventKey,
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
	list, err = tc.admin.ListAdjustmentTemplates(ctx, requestWithAdminAuth(&apiv1.ListAdjustmentTemplatesRequest{
		EventKey: eventKey,
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
	const eventKey = "test-event-key-admin-score"

	ctx := t.Context()
	tc := newTestClients()

	// Event started 45min ago — need stage IDs for score operations.
	ev := seedEvent(ctx, t, sharedTestDB, tc, eventKey, "NOW() + '-45 minutes'", 9)
	secondStageID := ev.stageIDs[1]

	p := seedPlayer(ctx, t, sharedTestDB, tc, "+15559380301", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE, "ScoreTestPlayer")

	// CreateStageScore via admin.
	createRes, err := tc.admin.CreateStageScore(ctx, requestWithAdminAuth(&apiv1.CreateStageScoreRequest{
		Data: &apiv1.StageScoreData{
			StageId:  secondStageID.String(),
			PlayerId: p.id.String(),
			Score: &apiv1.ScoreData{
				Value: 5,
			},
		},
	}))
	require.NoError(t, err, "CreateStageScore")
	assert.NotNil(t, createRes.Msg.GetScore(), "created score is returned")

	// Verify via GetScoresForPlayer.
	scores, err := tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err, "GetScoresForPlayer after create")

	scoreEntries := scores.Msg.GetScoreBoard().GetScores()
	require.NotEmpty(t, scoreEntries, "has score entries after create")

	// UpdateStageScore via admin.
	updateRes, err := tc.admin.UpdateStageScore(ctx, requestWithAdminAuth(&apiv1.UpdateStageScoreRequest{
		Score: &apiv1.StageScore{
			StageId:  secondStageID.String(),
			PlayerId: p.id.String(),
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

	// Verify updated score.
	scores, err = tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err, "GetScoresForPlayer after update")
	require.NotEmpty(t, scores.Msg.GetScoreBoard().GetScores(), "has score entries after update")

	// DeleteStageScore via admin.
	_, err = tc.admin.DeleteStageScore(ctx, requestWithAdminAuth(&apiv1.DeleteStageScoreRequest{
		StageId:  secondStageID.String(),
		PlayerId: p.id.String(),
	}))
	require.NoError(t, err, "DeleteStageScore")

	// Verify removal.
	listRes, err := tc.admin.ListStageScores(ctx, requestWithAdminAuth(&apiv1.ListStageScoresRequest{
		EventKey: eventKey,
	}))
	require.NoError(t, err, "ListStageScores after delete")

	for _, s := range listRes.Msg.GetScores() {
		assert.False(t,
			s.GetStageId() == secondStageID.String() && s.GetPlayerId() == p.id.String(),
			"deleted score should not appear in ListStageScores",
		)
	}
}

func Test_AdminPlayerManagement(t *testing.T) {
	const eventKey = "test-event-key-admin-player"

	ctx := t.Context()
	tc := newTestClients()

	// Insert bare event (no stages needed for player management).
	_, err := sharedTestDB.ExecContext(ctx, "INSERT INTO events (key, starts_at) VALUES ($1, NOW() + '1 day')", eventKey)
	require.NoError(t, err, "seed DB: insert future event")

	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	// CreatePlayer via admin with registration.
	p := seedPlayer(ctx, t, sharedTestDB, tc, "+15559380302", eventKey, apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE, "AdminCreated")

	// ListPlayers and verify the created player is present.
	list, err := tc.admin.ListPlayers(ctx, requestWithAdminAuth(&apiv1.ListPlayersRequest{
		EventKey: eventKey,
	}))
	require.NoError(t, err, "ListPlayers")

	var found bool

	for _, pl := range list.Msg.GetPlayers() {
		if pl.GetId() == p.id.String() {
			found = true

			assert.Equal(t, "AdminCreated", pl.GetData().GetName(), "player name matches")

			break
		}
	}

	assert.True(t, found, "created player found in ListPlayers")
}
