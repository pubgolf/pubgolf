package public

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao"
	"github.com/pubgolf/pubgolf/api/internal/lib/models"
	apiv1 "github.com/pubgolf/pubgolf/api/internal/lib/proto/api/v1"
	"github.com/pubgolf/pubgolf/api/internal/lib/sms"
)

func TestStartPlayerLogin(t *testing.T) {
	t.Parallel()

	playerID := models.PlayerIDFromULID(ulid.Make())
	phoneNum := models.PhoneNum("+15551234567")

	t.Run("New player", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{ID: playerID}, nil}}.Bind(mockDAO, "CreatePlayer")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(nil)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: "+15551234567",
		}))

		require.NoError(t, err)
	})

	t.Run("Existing player", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{}, dao.ErrAlreadyCreated}}.Bind(mockDAO, "CreatePlayer")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{false, nil}}.Bind(mockDAO, "PhoneNumberIsVerified")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(nil)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: "+15551234567",
		}))

		require.NoError(t, err)
	})

	t.Run("Invalid phone number", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: "invalid",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("SMS send failure", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{ID: playerID}, nil}}.Bind(mockDAO, "CreatePlayer")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(sms.ErrUpstreamProviderIssue)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: "+15551234567",
		}))

		require.Error(t, err)
	})
}
