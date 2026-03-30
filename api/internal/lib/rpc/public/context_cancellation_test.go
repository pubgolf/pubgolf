package public

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

// mockDAOCallWithError sets up a DAO mock method to return the given error with
// a zero-value result. argCount specifies the number of arguments (all matched
// with mock.Anything).
func mockDAOCallWithError(m *dao.MockQueryProvider, method string, zeroVal any, ctxErr error, argCount int) {
	args := make([]any, argCount)
	for i := range args {
		args[i] = mock.Anything
	}

	dao.MockDAOCall{
		ShouldCall: true,
		Args:       args,
		Return:     []any{zeroVal, ctxErr},
	}.Bind(m, method)
}

// contextErrorCases returns table-driven test cases for both flavors of context
// error: a pre-canceled context and an already-expired deadline.
func contextErrorCases() []struct {
	name    string
	makeCtx func(t *testing.T) context.Context
	ctxErr  error
} {
	return []struct {
		name    string
		makeCtx func(t *testing.T) context.Context
		ctxErr  error
	}{
		{
			name: "canceled context",
			makeCtx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()

				return ctx
			},
			ctxErr: context.Canceled,
		},
		{
			name: "expired deadline",
			makeCtx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
				t.Cleanup(cancel)

				return ctx
			},
			ctxErr: context.DeadlineExceeded,
		},
	}
}

func TestGuardRegisteredForEvent_ContextErrors(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())
	eventKey := "test-event"

	for _, ce := range contextErrorCases() {
		t.Run(ce.name, func(t *testing.T) {
			t.Parallel()

			t.Run("propagates from EventIDByKey", func(t *testing.T) {
				t.Parallel()

				m := new(dao.MockQueryProvider)
				mockDAOCallWithError(m, "EventIDByKey", models.EventID{}, ce.ctxErr, 2)
				s := makeTestServer(m)

				_, err := s.guardRegisteredForEvent(ce.makeCtx(t), playerID, eventKey)

				require.Error(t, err)
				assert.ErrorIs(t, err, ce.ctxErr)
			})

			t.Run("propagates from PlayerRegisteredForEvent mid-request", func(t *testing.T) {
				t.Parallel()

				m := new(dao.MockQueryProvider)
				mockEventIDByKey(m, eventKey, eventID)
				mockDAOCallWithError(m, "PlayerRegisteredForEvent", false, ce.ctxErr, 3)
				s := makeTestServer(m)

				_, err := s.guardRegisteredForEvent(ce.makeCtx(t), playerID, eventKey)

				require.Error(t, err)
				assert.ErrorIs(t, err, ce.ctxErr)
			})
		})
	}
}

func TestGuardStageID_ContextErrors(t *testing.T) {
	t.Parallel()

	eventID := models.EventIDFromULID(ulid.Make())
	venueKey := models.VenueKeyFromUInt32(1)

	for _, ce := range contextErrorCases() {
		t.Run(ce.name, func(t *testing.T) {
			t.Parallel()

			m := new(dao.MockQueryProvider)
			mockDAOCallWithError(m, "StageIDByVenueKey", models.StageID{}, ce.ctxErr, 3)
			s := makeTestServer(m)

			_, err := s.guardStageID(ce.makeCtx(t), eventID, venueKey)

			require.Error(t, err)
			assert.ErrorIs(t, err, ce.ctxErr)
		})
	}
}

func TestGuardPlayerCategory_ContextErrors(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	eventID := models.EventIDFromULID(ulid.Make())

	for _, ce := range contextErrorCases() {
		t.Run(ce.name, func(t *testing.T) {
			t.Parallel()

			m := new(dao.MockQueryProvider)
			mockDAOCallWithError(m, "PlayerCategoryForEvent", models.ScoringCategoryUnspecified, ce.ctxErr, 3)
			s := makeTestServer(m)

			_, err := s.guardPlayerCategory(ce.makeCtx(t), playerID, eventID)

			require.Error(t, err)
			assert.ErrorIs(t, err, ce.ctxErr)
		})
	}
}
