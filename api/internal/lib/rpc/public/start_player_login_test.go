package public

import (
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

func TestStartPlayerLogin(t *testing.T) {
	t.Parallel()

	phoneNum := models.PhoneNum("+1" + faker.E164PhoneNumber()[2:])
	playerID := models.PlayerIDFromULID(ulid.Make())

	t.Run("Sends verification to new player", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{ID: playerID}, nil}}.Bind(mockDAO, "CreatePlayer")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(nil)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: string(phoneNum),
		}))

		require.NoError(t, err)
	})

	t.Run("Sends verification to existing player", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{}, dao.ErrAlreadyCreated}}.Bind(mockDAO, "CreatePlayer")
		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, phoneNum}, Return: []any{false, nil}}.Bind(mockDAO, "PhoneNumberIsVerified")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(nil)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: string(phoneNum),
		}))

		require.NoError(t, err)
	})

	t.Run("Rejects invalid phone numbers", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name  string
			phone string
		}{
			{name: "empty", phone: ""},
			{name: "letters", phone: "invalid"},
			{name: "missing plus", phone: "12345678901"},
			{name: "alpha after plus", phone: "+1abc"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				mockDAO := new(dao.MockQueryProvider)
				mockMes := new(sms.MockMessenger)
				s := NewServer(mockDAO, mockMes)

				_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
					PhoneNumber: tc.phone,
				}))

				require.Error(t, err)
				assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
			})
		}
	})

	t.Run("Returns error on SMS send failure", func(t *testing.T) {
		t.Parallel()

		mockDAO := new(dao.MockQueryProvider)
		mockMes := new(sms.MockMessenger)
		s := NewServer(mockDAO, mockMes)

		dao.MockDAOCall{ShouldCall: true, Args: []any{mock.Anything, "", phoneNum}, Return: []any{models.Player{ID: playerID}, nil}}.Bind(mockDAO, "CreatePlayer")
		mockMes.On("SendVerification", mock.Anything, phoneNum).Return(sms.ErrUpstreamProviderIssue)

		_, err := s.StartPlayerLogin(t.Context(), connect.NewRequest(&apiv1.StartPlayerLoginRequest{
			PhoneNumber: string(phoneNum),
		}))

		require.Error(t, err)
		assert.Equal(t, connect.CodeUnavailable, connect.CodeOf(err))
	})
}
