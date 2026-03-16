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

func TestCompletePlayerLogin(t *testing.T) {
	t.Parallel()

	t.Run("Valid login", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		phoneNum := models.PhoneNum("+15551234567")
		playerID := models.PlayerIDFromULID(ulid.Make())
		authToken := models.AuthTokenFromULID(ulid.Make())

		mockMes.On("CheckVerification", mock.Anything, phoneNum, "12345678").Return(true, nil)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{true, nil}}.Bind(mockDAO, "VerifyPhoneNumber")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{dao.GenerateAuthTokenResult{PlayerID: playerID, AuthToken: authToken}, nil}}.Bind(mockDAO, "GenerateAuthToken")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID}, Return: []any{models.Player{ID: playerID, Name: "Test"}, nil}}.Bind(mockDAO, "PlayerByID")

		resp, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: "+15551234567",
			AuthCode:    "12345678",
		}))

		require.NoError(t, err)
		assert.Equal(t, authToken.String(), resp.Msg.GetAuthToken())
		assert.Equal(t, "Test", resp.Msg.GetPlayer().GetData().GetName())
	})

	t.Run("Invalid phone format", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: "invalid",
			AuthCode:    "12345678",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("Invalid auth code format", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: "+15551234567",
			AuthCode:    "abc",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
	})

	t.Run("SMS check fails (wrong code)", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		phoneNum := models.PhoneNum("+15551234567")

		mockMes.On("CheckVerification", mock.Anything, phoneNum, "12345678").Return(false, nil)

		_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: "+15551234567",
			AuthCode:    "12345678",
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})
}
