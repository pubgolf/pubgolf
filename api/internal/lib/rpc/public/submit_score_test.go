package public

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/forms"
	"github.com/pubgolf/pubgolf/api/internal/lib/middleware"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
)

func TestSubmitScore(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	gameCtx := middleware.ContextWithPlayerID(t.Context(), playerID)
	eventKey := "test-event"
	eventID := models.EventIDFromULID(ulid.Make())
	stageID := models.StageIDFromULID(ulid.Make())
	adjTemplateID := models.AdjustmentTemplateIDFromULID(ulid.Make())

	mockGuards := func(mockDAO *dao.MockQueryProvider) {
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{true, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventID, models.VenueKeyFromUInt32(1)}, Return: []any{stageID, nil}}.Bind(mockDAO, "StageIDByVenueKey")
	}

	validFormData := &apiv1.FormSubmission{
		Values: []*apiv1.FormValue{
			{
				Id: forms.SubmitScoreInputIDSips,
				Value: &apiv1.FormValue_Numeric{
					Numeric: 3,
				},
			},
		},
	}

	t.Run("Valid score, no adjustments", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockGuards(mockDAO)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, stageID}, Return: []any{[]models.AdjustmentTemplate(nil), nil}}.Bind(mockDAO, "AdjustmentTemplatesByStageID")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}, Return: []any{nil}}.Bind(mockDAO, "UpsertScore")

		resp, err := s.SubmitScore(gameCtx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
			Data:     validFormData,
		}))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())
	})

	t.Run("Valid score with adjustments", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockGuards(mockDAO)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, stageID}, Return: []any{[]models.AdjustmentTemplate{
			{ID: adjTemplateID, Label: "Bonus", Value: -1, VenueSpecific: false},
		}, nil}}.Bind(mockDAO, "AdjustmentTemplatesByStageID")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}, Return: []any{nil}}.Bind(mockDAO, "UpsertScore")

		resp, err := s.SubmitScore(gameCtx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: forms.SubmitScoreInputIDSips,
						Value: &apiv1.FormValue_Numeric{
							Numeric: 3,
						},
					},
					{
						Id: forms.SubmitScoreInputIDStandardAdj,
						Value: &apiv1.FormValue_SelectMany{
							SelectMany: &apiv1.SelectManyValue{
								SelectedIds: []string{adjTemplateID.String()},
							},
						},
					},
				},
			},
		}))

		require.NoError(t, err)
		assert.Equal(t, apiv1.ScoreStatus_SCORE_STATUS_SUBMITTED_EDITABLE, resp.Msg.GetStatus())
	})

	t.Run("Invalid form data (score=0)", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockGuards(mockDAO)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, stageID}, Return: []any{[]models.AdjustmentTemplate(nil), nil}}.Bind(mockDAO, "AdjustmentTemplatesByStageID")

		_, err := s.SubmitScore(gameCtx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: forms.SubmitScoreInputIDSips,
						Value: &apiv1.FormValue_Numeric{
							Numeric: 0,
						},
					},
				},
			},
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("Unknown adjustment ID", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		mockGuards(mockDAO)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, stageID}, Return: []any{[]models.AdjustmentTemplate(nil), nil}}.Bind(mockDAO, "AdjustmentTemplatesByStageID")

		unknownAdjID := models.AdjustmentTemplateIDFromULID(ulid.Make())

		_, err := s.SubmitScore(gameCtx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
			Data: &apiv1.FormSubmission{
				Values: []*apiv1.FormValue{
					{
						Id: forms.SubmitScoreInputIDSips,
						Value: &apiv1.FormValue_Numeric{
							Numeric: 3,
						},
					},
					{
						Id: forms.SubmitScoreInputIDStandardAdj,
						Value: &apiv1.FormValue_SelectMany{
							SelectMany: &apiv1.SelectManyValue{
								SelectedIds: []string{unknownAdjID.String()},
							},
						},
					},
				},
			},
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("Not registered for event", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		s := makeTestServer(mockDAO)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, eventKey}, Return: []any{eventID, nil}}.Bind(mockDAO, "EventIDByKey")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID, eventID}, Return: []any{false, nil}}.Bind(mockDAO, "PlayerRegisteredForEvent")

		_, err := s.SubmitScore(gameCtx, connect.NewRequest(&apiv1.SubmitScoreRequest{
			PlayerId: playerID.String(),
			EventKey: eventKey,
			VenueKey: 1,
			Data:     validFormData,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})
}
