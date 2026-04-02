package public

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/go-faker/faker/v4"
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

	phoneNum := models.PhoneNum("+1" + faker.E164PhoneNumber()[2:])
	authCode := sms.MockAuthCode

	t.Run("Returns player and token on valid login", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		playerID := models.PlayerIDFromULID(ulid.Make())
		authToken := models.AuthTokenFromULID(ulid.Make())
		playerName := faker.Name()

		mockMes.On("CheckVerification", mock.Anything, phoneNum, authCode).Return(true, nil)
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{true, nil}}.Bind(mockDAO, "VerifyPhoneNumber")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{dao.GenerateAuthTokenResult{PlayerID: playerID, AuthToken: authToken}, nil}}.Bind(mockDAO, "GenerateAuthToken")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, playerID}, Return: []any{models.Player{ID: playerID, Name: playerName}, nil}}.Bind(mockDAO, "PlayerByID")

		resp, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: string(phoneNum),
			AuthCode:    authCode,
		}))

		require.NoError(t, err)
		assert.Equal(t, authToken.String(), resp.Msg.GetAuthToken())
		assert.Equal(t, playerName, resp.Msg.GetPlayer().GetData().GetName())
	})

	t.Run("Rejects invalid inputs", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name  string
			phone string
			code  string
		}{
			{name: "empty phone", phone: "", code: authCode},
			{name: "letters-only phone", phone: "invalid", code: authCode},
			{name: "phone without plus", phone: "15551234567", code: authCode},
			{name: "empty auth code", phone: string(phoneNum), code: ""},
			{name: "short auth code", phone: string(phoneNum), code: "1234567"},
			{name: "auth code with letters", phone: string(phoneNum), code: "1234abcd"},
			{name: "alpha auth code", phone: string(phoneNum), code: "abcdefgh"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				mockMes := new(sms.MockMessenger)
				s := NewServer(mockDAO, mockMes)

				_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
					PhoneNumber: tc.phone,
					AuthCode:    tc.code,
				}))

				require.Error(t, err)
				assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
			})
		}
	})

	t.Run("Rejects wrong verification code", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		mockMes.On("CheckVerification", mock.Anything, phoneNum, authCode).Return(false, nil)

		_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: string(phoneNum),
			AuthCode:    authCode,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))
	})

	t.Run("Returns error on SMS check failure", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		mockMes.On("CheckVerification", mock.Anything, phoneNum, authCode).Return(false, errors.New("twilio down"))

		_, err := s.CompletePlayerLogin(t.Context(), connect.NewRequest(&apiv1.CompletePlayerLoginRequest{
			PhoneNumber: string(phoneNum),
			AuthCode:    authCode,
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})
}
