// admin_test.go — admin management flow: CRUD operations for players, stages,
// venues, scores, and adjustment templates.

//nolint:paralleltest
package e2e

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func Test_AdminAdjustmentTemplates(t *testing.T) {
	const eventKey = "test-event-key-admin-adj"

	ctx := t.Context()
	tc := newTestClients()

	seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '-45 minutes'",
		NumStages:    9,
	})

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
	ev := seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '-45 minutes'",
		NumStages:    9,
	})
	secondStageID := ev.stageIDs[1]

	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380301",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		Name:     "ScoreTestPlayer",
	})

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

	// Verify score value via GetScoresForPlayer.
	scores, err := tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err, "GetScoresForPlayer after create")

	scoreEntries := scores.Msg.GetScoreBoard().GetScores()
	require.NotEmpty(t, scoreEntries, "has score entries after create")

	hasScore5 := slices.ContainsFunc(scoreEntries, func(e *apiv1.ScoreBoard_ScoreBoardEntry) bool {
		return e.GetScore() == 5
	})
	assert.True(t, hasScore5, "score of 5 appears in player scoreboard after create")

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

	// Verify updated score value.
	scores, err = tc.pub.GetScoresForPlayer(ctx, requestWithAuth(&apiv1.GetScoresForPlayerRequest{
		EventKey: eventKey,
		PlayerId: p.id.String(),
	}, p.token))
	require.NoError(t, err, "GetScoresForPlayer after update")

	scoreEntries = scores.Msg.GetScoreBoard().GetScores()
	require.NotEmpty(t, scoreEntries, "has score entries after update")

	hasScore7 := slices.ContainsFunc(scoreEntries, func(e *apiv1.ScoreBoard_ScoreBoardEntry) bool {
		return e.GetScore() == 7
	})
	assert.True(t, hasScore7, "score of 7 appears in player scoreboard after update")

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

	seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '1 day'",
	})

	// CreatePlayer via admin with registration.
	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380302",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
		Name:     "AdminCreated",
	})

	// ListPlayers and verify the created player is present.
	list, err := tc.admin.ListPlayers(ctx, requestWithAdminAuth(&apiv1.ListPlayersRequest{
		EventKey: eventKey,
	}))
	require.NoError(t, err, "ListPlayers")

	playerInList := slices.ContainsFunc(list.Msg.GetPlayers(), func(pl *apiv1.Player) bool {
		return pl.GetId() == p.id.String() && pl.GetData().GetName() == "AdminCreated"
	})
	assert.True(t, playerInList, "created player found in ListPlayers with correct name")

	// Update the player's name and category via admin.
	updatedName := "AdminUpdated"
	updatedCategory := apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_FIVE_HOLE

	updateRes, err := tc.admin.UpdatePlayer(ctx, requestWithAdminAuth(&apiv1.UpdatePlayerRequest{
		PlayerId: p.id.String(),
		PlayerData: &apiv1.PlayerData{
			Name: updatedName,
		},
		Registration: &apiv1.EventRegistration{
			EventKey:        eventKey,
			ScoringCategory: updatedCategory,
		},
	}))
	require.NoError(t, err, "UpdatePlayer")
	assert.Equal(t, updatedName, updateRes.Msg.GetPlayer().GetData().GetName(), "updated name returned")
}

func Test_AdminEventSetup(t *testing.T) {
	const eventKey = "test-event-key-admin-event"

	ctx := t.Context()
	tc := newTestClients()

	ev := seedEvent(ctx, t, sharedTestDB, tc, seedEventOpts{
		EventKey:     eventKey,
		StartsAtExpr: "NOW() + '-45 minutes'",
		NumStages:    3,
	})

	// Register a player so we can verify schedule propagation after stage updates.
	p := seedPlayer(ctx, t, sharedTestDB, tc, seedPlayerOpts{
		Phone:    "+15559380401",
		EventKey: eventKey,
		Category: apiv1.ScoringCategory_SCORING_CATEGORY_PUB_GOLF_NINE_HOLE,
	})

	// ListVenues — verify seeded venues are visible.
	venueList, err := tc.admin.ListVenues(ctx, requestWithAdminAuth(&apiv1.ListVenuesRequest{}))
	require.NoError(t, err, "ListVenues")
	assert.GreaterOrEqual(t, len(venueList.Msg.GetVenues()), 3, "at least 3 venues exist")

	// ListEventStages — capture existing stage state.
	stageList, err := tc.admin.ListEventStages(ctx, requestWithAdminAuth(&apiv1.ListEventStagesRequest{
		EventKey: eventKey,
	}))
	require.NoError(t, err, "ListEventStages")
	require.Len(t, stageList.Msg.GetStages(), 3, "3 stages in event")

	// Use the second stage (current stage at -45min) for the update test,
	// so we can verify the venue key change propagates to GetSchedule's currentVenueKey.
	currentStage := stageList.Msg.GetStages()[1]
	originalDuration := currentStage.GetDurationMin()
	currentVenueID := currentStage.GetVenue().GetId()

	// Get schedule before update to capture venue key (a cache key that changes on stage updates).
	scheduleBefore, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule before UpdateStage")

	venueKeyBefore := scheduleBefore.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotZero(t, venueKeyBefore, "has current venue key before update")

	visitedKeysBefore := scheduleBefore.Msg.GetSchedule().GetVisitedVenueKeys()
	require.NotEmpty(t, visitedKeysBefore, "has visited venue keys before update")

	// UpdateStage — change current stage's duration to a different value.
	// This also sets venue_key = NULL, which forces reassignment on next schedule fetch.
	newDuration := originalDuration + 15

	_, err = tc.admin.UpdateStage(ctx, requestWithAdminAuth(&apiv1.UpdateStageRequest{
		StageId:     ev.stageIDs[1].String(),
		VenueId:     currentVenueID,
		Rank:        currentStage.GetRank(),
		DurationMin: newDuration,
	}))
	require.NoError(t, err, "UpdateStage")

	// Verify duration updated via ListEventStages.
	stageList, err = tc.admin.ListEventStages(ctx, requestWithAdminAuth(&apiv1.ListEventStagesRequest{
		EventKey: eventKey,
	}))
	require.NoError(t, err, "ListEventStages after update")

	updatedStage := stageList.Msg.GetStages()[1]
	assert.Equal(t, newDuration, updatedStage.GetDurationMin(), "stage duration updated")
	assert.NotEqual(t, originalDuration, updatedStage.GetDurationMin(), "duration changed from original")

	// Verify venue key changed (cache invalidation propagated to schedule).
	_, err = tc.admin.PurgeAllCaches(ctx, requestWithAdminAuth(&apiv1.PurgeAllCachesRequest{}))
	require.NoError(t, err)

	scheduleAfter, err := tc.pub.GetSchedule(ctx, requestWithAuth(&apiv1.GetScheduleRequest{
		EventKey: eventKey,
	}, p.token))
	require.NoError(t, err, "GetSchedule after UpdateStage")

	venueKeyAfter := scheduleAfter.Msg.GetSchedule().GetCurrentVenueKey()
	require.NotZero(t, venueKeyAfter, "has current venue key after update")
	assert.NotEqual(t, venueKeyBefore, venueKeyAfter, "updated stage's venue key changed (cache invalidation)")

	// Non-updated stages should retain their venue keys.
	visitedKeysAfter := scheduleAfter.Msg.GetSchedule().GetVisitedVenueKeys()
	assert.Equal(t, visitedKeysBefore, visitedKeysAfter, "non-updated stage venue keys unchanged")
}
