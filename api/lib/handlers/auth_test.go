package handlers_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "github.com/escavelo/pubgolf/api/lib/handlers"
	"github.com/escavelo/pubgolf/api/lib/utils"
	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

const (
	emptyEventID                  = "00000000-0000-4000-a004-000000000000"
	emptyEventKey                 = "empty-event"
	eventWithRegistrationsID      = "00000000-0000-4000-a004-000000000001"
	eventWithRegistrationsKey     = "event-with-registrations"
	freshRegistrationPhoneNumber  = "+15550000000"
	freshRegistrationAuthCode     = 123456
	loggedOutConfirmedPhoneNumber = "+15550000001"
	loggedInConfirmedPhoneNumber  = "+15550000002"
)

func queryPropertyFromUserRow(propertyToSelect string, playerPhoneNumber string, eventID string) string {
	return fmt.Sprintf(`
		SELECT %s
		FROM players
		WHERE phone_number = '%s'
			AND event_id = '%s';
	`, propertyToSelect, playerPhoneNumber, eventID)
}

func TestMain(m *testing.M) {
	utils.TestMain(m, &testDB, []string{
		fmt.Sprintf(`INSERT INTO events(
			id,
			key,
			name,
			start_time,
			end_time
		) VALUES (
			'%s',
			'%s',
			'Empty Event',
			NOW() + interval '30m',
			NOW() + interval '6h30m'
		);`, emptyEventID, emptyEventKey),

		fmt.Sprintf(`INSERT INTO events(
			id,
			key,
			name,
			start_time,
			end_time
			) VALUES (
			'%s',
			'%s',
			'Event with Registrations',
			NOW() + interval '30m',
			NOW() + interval '6h30m'
			);`, eventWithRegistrationsID, eventWithRegistrationsKey),
		fmt.Sprintf(`INSERT INTO players(
			event_id,
			name,
			phone_number,
			league,
			phone_number_confirmed,
			auth_code,
			auth_code_created_at,
			auth_token
		) VALUES (
			'%s',
			'Fresh Registration',
			'%s',
			'PGA',
			FALSE,
			%d,
			NOW(),
			NULL
		);`, eventWithRegistrationsID, freshRegistrationPhoneNumber, freshRegistrationAuthCode),
		fmt.Sprintf(`INSERT INTO players(
			event_id,
			name,
			phone_number,
			league,
			phone_number_confirmed,
			auth_code,
			auth_code_created_at,
			auth_token
		) VALUES (
			'%s',
			'Logged Out with Confirmed Number',
			'%s',
			'PGA',
			TRUE,
			123456,
			NOW(),
			NULL
		);`, eventWithRegistrationsID, loggedOutConfirmedPhoneNumber),
		fmt.Sprintf(`INSERT INTO players(
			event_id,
			name,
			phone_number,
			league,
			phone_number_confirmed,
			auth_code,
			auth_code_created_at,
			auth_token
		) VALUES (
			'%s',
			'Logged In with Confirmed Number',
			'%s',
			'PGA',
			TRUE,
			NULL,
			NOW(),
			uuid_generate_v4 ()
		);`, eventWithRegistrationsID, loggedInConfirmedPhoneNumber),
	})
}

func TestRegisterPlayer(t *testing.T) {
	testCases := []struct {
		testDescription   string
		providedReq       *pg.RegisterPlayerRequest
		expectedRep       *pg.RegisterPlayerReply
		expectedReplyCode codes.Code
	}{
		{
			"Missing event key in request returns InvalidArgument",
			&pg.RegisterPlayerRequest{
				PhoneNumber: freshRegistrationPhoneNumber,
				Name:        "Bob Loblaw",
			},
			&pg.RegisterPlayerReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid event key in request returns NotFound",
			&pg.RegisterPlayerRequest{
				EventKey:    "fake-event-key",
				PhoneNumber: freshRegistrationPhoneNumber,
				Name:        "Bob Loblaw",
			},
			&pg.RegisterPlayerReply{},
			codes.NotFound,
		},
		{
			"Missing phone number in request returns InvalidArgument",
			&pg.RegisterPlayerRequest{
				EventKey: "fake-event-key",
				Name:     "Bob Loblaw",
			},
			&pg.RegisterPlayerReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid phone number in request returns InvalidArgument",
			&pg.RegisterPlayerRequest{
				EventKey:    "fake-event-key",
				PhoneNumber: "15550000000",
				Name:        "Bob Loblaw",
			},
			&pg.RegisterPlayerReply{},
			codes.InvalidArgument,
		},
		{
			"Missing name in request returns InvalidArgument",
			&pg.RegisterPlayerRequest{
				EventKey:    "fake-event-key",
				PhoneNumber: freshRegistrationPhoneNumber,
			},
			&pg.RegisterPlayerReply{},
			codes.InvalidArgument,
		},
		{
			"Valid request data returns OK and empty response",
			&pg.RegisterPlayerRequest{
				EventKey:    emptyEventKey,
				PhoneNumber: freshRegistrationPhoneNumber,
				Name:        "Bob Loblaw",
			},
			&pg.RegisterPlayerReply{},
			codes.OK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(tt *testing.T) {
			rd, _ := makeTestUnauthenticatedRequestData()
			defer rd.Tx.Rollback()

			rep, err := RegisterPlayer(rd, tc.providedReq)

			// Only compare the response if there isn't an error expected.
			if tc.expectedReplyCode == codes.OK {
				assert.NoError(tt, err)
				assert.Equal(tt, tc.expectedRep, rep)
			}

			st, ok := status.FromError(err)
			assert.True(tt, ok, "Unable to parse status code from error")
			assert.Equal(tt, tc.expectedReplyCode.String(), st.Code().String())
		})
	}
}

func TestRegisterPlayerRejectsDuplicatePhoneNumbers(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	_, err := RegisterPlayer(rd, &pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	})
	assert.NoError(t, err)

	_, err = RegisterPlayer(rd, &pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw's Brother",
	})
	st, ok := status.FromError(err)
	assert.True(t, ok, "Unable to parse status code from error")
	assert.Equal(t, codes.AlreadyExists.String(), st.Code().String())
}

func TestRegisterPlayerAcceptsDuplicatePhoneNumbersAcrossEvents(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	_, err := RegisterPlayer(rd, &pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	})
	assert.NoError(t, err)

	_, err = rd.Tx.Exec(`
		INSERT INTO events(
			key,
			name,
			start_time,
			end_time
		) VALUES (
			'empty-event-2',
			'Empty Event 2',
			NOW() + interval '30m',
			NOW() + interval '6h30m'
		);
	`)
	assert.NoError(t, err)

	_, err = RegisterPlayer(rd, &pg.RegisterPlayerRequest{
		EventKey:    "empty-event-2",
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	})
	assert.NoError(t, err)
}

func TestRegisterPlayerCreatesRowInDB(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("COUNT(*)", freshRegistrationPhoneNumber, emptyEventID)

	var playerCount int
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&playerCount)
	assert.NoError(t, err)
	assert.Equal(t, 0, playerCount)

	req := pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	}
	_, err = RegisterPlayer(rd, &req)
	assert.NoError(t, err)

	row = rd.Tx.QueryRow(query)
	err = row.Scan(&playerCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, playerCount)
}

func TestRegisterPlayerGeneratesAuthCode(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_code", freshRegistrationPhoneNumber, emptyEventID)

	var authCode sql.NullInt32
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authCode)
	assert.Error(t, err) // sql: no rows in result set
	assert.False(t, authCode.Valid)

	req := pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	}
	_, err = RegisterPlayer(rd, &req)
	assert.NoError(t, err)

	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authCode)
	assert.NoError(t, err)
	assert.True(t, authCode.Valid)
}

func TestRegisterPlayerSendsSms(t *testing.T) {
	rd, logHook := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	req := pg.RegisterPlayerRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		Name:        "Bob Loblaw",
	}
	_, err := RegisterPlayer(rd, &req)
	assert.NoError(t, err)

	assert.Equal(t, 1, countLogEntries(logHook, "Twilio Simulation"))
}

func TestRequestPlayerLogin(t *testing.T) {
	testCases := []struct {
		testDescription   string
		providedReq       *pg.RequestPlayerLoginRequest
		expectedRep       *pg.RequestPlayerLoginReply
		expectedReplyCode codes.Code
	}{
		{
			"Missing event key in request returns InvalidArgument",
			&pg.RequestPlayerLoginRequest{
				PhoneNumber: freshRegistrationPhoneNumber,
			},
			&pg.RequestPlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid event key in request returns NotFound",
			&pg.RequestPlayerLoginRequest{
				EventKey:    "fake-event-key",
				PhoneNumber: freshRegistrationPhoneNumber,
			},
			&pg.RequestPlayerLoginReply{},
			codes.NotFound,
		},
		{
			"Missing phone number in request returns InvalidArgument",
			&pg.RequestPlayerLoginRequest{
				EventKey: eventWithRegistrationsKey,
			},
			&pg.RequestPlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid phone number in request returns InvalidArgument",
			&pg.RequestPlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: "15550000000",
			},
			&pg.RequestPlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Valid request data returns OK and empty response",
			&pg.RequestPlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: freshRegistrationPhoneNumber,
			},
			&pg.RequestPlayerLoginReply{},
			codes.OK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(tt *testing.T) {
			rd, _ := makeTestUnauthenticatedRequestData()
			defer rd.Tx.Rollback()

			rep, err := RequestPlayerLogin(rd, tc.providedReq)

			// Only compare the response if there isn't an error expected.
			if tc.expectedReplyCode == codes.OK {
				assert.NoError(tt, err)
				assert.Equal(tt, tc.expectedRep, rep)
			}

			st, ok := status.FromError(err)
			assert.True(tt, ok, "Unable to parse status code from error")
			assert.Equal(tt, tc.expectedReplyCode.String(), st.Code().String())
		})
	}
}

func TestRequestPlayerLoginReturnsOKWithSmsForRegisteredNumber(t *testing.T) {
	rd, logHook := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	_, err := RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: freshRegistrationPhoneNumber,
	})
	assert.NoError(t, err)

	assert.Equal(t, 1, countLogEntries(logHook, "Twilio Simulation"))
}

func TestRequestPlayerLoginReturnsOKWithoutSmsForUnregisteredNumber(t *testing.T) {
	rd, logHook := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	_, err := RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: "+15555555555",
	})
	assert.NoError(t, err)

	assert.Equal(t, 0, countLogEntries(logHook, "Twilio Simulation"))
}

func TestRequestPlayerLoginReturnsOKWithoutSmsForWrongEventKey(t *testing.T) {
	rd, logHook := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	_, err := RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    emptyEventKey,
		PhoneNumber: freshRegistrationPhoneNumber,
	})
	assert.NoError(t, err)

	assert.Equal(t, 0, countLogEntries(logHook, "Twilio Simulation"))
}

func TestRequestPlayerLoginRegeneratesAuthCode(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_code", loggedOutConfirmedPhoneNumber, eventWithRegistrationsID)

	var authCodeOld int
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authCodeOld)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, authCodeOld)

	_, err = RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: loggedOutConfirmedPhoneNumber,
	})
	assert.NoError(t, err)

	var authCodeNew int
	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authCodeNew)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, authCodeNew)

	assert.NotEqual(t, authCodeOld, authCodeNew)
}

func TestRequestPlayerLoginRegeneratesAuthCodeExpiration(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_code_created_at", loggedOutConfirmedPhoneNumber, eventWithRegistrationsID)

	var authCodeTimestampOld time.Time
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authCodeTimestampOld)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, authCodeTimestampOld)

	_, err = RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: loggedOutConfirmedPhoneNumber,
	})
	assert.NoError(t, err)

	var authCodeTimestampNew time.Time
	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authCodeTimestampNew)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, authCodeTimestampNew)

	assert.True(t, authCodeTimestampOld.Before(authCodeTimestampNew))
}

func TestRequestPlayerLoginInvalidatesAuthToken(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_token", loggedInConfirmedPhoneNumber, eventWithRegistrationsID)

	var authToken sql.NullString
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authToken)
	assert.NoError(t, err)
	assert.True(t, authToken.Valid)

	// User is "Logged In with Confirmed Number".
	_, err = RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: loggedInConfirmedPhoneNumber,
	})
	assert.NoError(t, err)

	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authToken)
	assert.NoError(t, err)
	assert.False(t, authToken.Valid)
}

func TestRequestPlayerLoginResetsUpdatedAtTimestampInUserRow(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("updated_at", freshRegistrationPhoneNumber, eventWithRegistrationsID)

	var userEditTimestampOld time.Time
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&userEditTimestampOld)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userEditTimestampOld)

	// User is "Logged Out with Confirmed Number".
	_, err = RequestPlayerLogin(rd, &pg.RequestPlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: freshRegistrationPhoneNumber,
	})
	assert.NoError(t, err)

	var userEditTimestampNew time.Time
	row = rd.Tx.QueryRow(query)
	err = row.Scan(&userEditTimestampNew)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userEditTimestampNew)

	assert.True(t, userEditTimestampOld.Before(userEditTimestampNew))
}

func TestPlayerLogin(t *testing.T) {
	testCases := []struct {
		testDescription   string
		providedReq       *pg.PlayerLoginRequest
		expectedRep       *pg.PlayerLoginReply
		expectedReplyCode codes.Code
	}{
		{
			"Missing event key in request returns InvalidArgument",
			&pg.PlayerLoginRequest{
				PhoneNumber: freshRegistrationPhoneNumber,
				AuthCode:    freshRegistrationAuthCode,
			},
			&pg.PlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid event key in request returns NotFound",
			&pg.PlayerLoginRequest{
				EventKey:    "fake-event-key",
				PhoneNumber: freshRegistrationPhoneNumber,
				AuthCode:    freshRegistrationAuthCode,
			},
			&pg.PlayerLoginReply{},
			codes.NotFound,
		},
		{
			"Missing phone number in request returns InvalidArgument",
			&pg.PlayerLoginRequest{
				EventKey: eventWithRegistrationsKey,
				AuthCode: freshRegistrationAuthCode,
			},
			&pg.PlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid phone number in request returns InvalidArgument",
			&pg.PlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: "15550000000",
				AuthCode:    freshRegistrationAuthCode,
			},
			&pg.PlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Missing auth code in request returns InvalidArgument",
			&pg.PlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: freshRegistrationPhoneNumber,
			},
			&pg.PlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Invalid auth code in request returns InvalidArgument",
			&pg.PlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: freshRegistrationPhoneNumber,
				AuthCode:    42,
			},
			&pg.PlayerLoginReply{},
			codes.InvalidArgument,
		},
		{
			"Incorrect auth code returns Unauthenticated",
			&pg.PlayerLoginRequest{
				EventKey:    eventWithRegistrationsKey,
				PhoneNumber: freshRegistrationPhoneNumber,
				AuthCode:    987654,
			},
			&pg.PlayerLoginReply{},
			codes.Unauthenticated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(tt *testing.T) {
			rd, _ := makeTestUnauthenticatedRequestData()
			defer rd.Tx.Rollback()

			rep, err := PlayerLogin(rd, tc.providedReq)

			// Only compare the response if there isn't an error expected.
			if tc.expectedReplyCode == codes.OK {
				assert.NoError(tt, err)
				assert.Equal(tt, tc.expectedRep, rep)
			}

			st, ok := status.FromError(err)
			assert.True(tt, ok, "Unable to parse status code from error")
			assert.Equal(tt, tc.expectedReplyCode.String(), st.Code().String())
		})
	}
}

func TestPlayerLoginReturnsUserData(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	rep, err := PlayerLogin(rd, &pg.PlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		AuthCode:    freshRegistrationAuthCode,
	})
	assert.NoError(t, err)

	assert.NotEmpty(t, rep.AuthToken)
	assert.NotEmpty(t, rep.PlayerID)
	assert.Equal(t, pg.PlayerRole_DEFAULT, rep.GetPlayerRole())
}

func TestPlayerLoginUnSetsAuthCode(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_code", freshRegistrationPhoneNumber, eventWithRegistrationsID)

	var authCode sql.NullInt32
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authCode)
	assert.NoError(t, err)
	assert.True(t, authCode.Valid)

	_, err = PlayerLogin(rd, &pg.PlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		AuthCode:    freshRegistrationAuthCode,
	})
	assert.NoError(t, err)

	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authCode)
	assert.NoError(t, err)
	assert.False(t, authCode.Valid)
}

func TestPlayerLoginPersistsAuthTokenToDB(t *testing.T) {
	rd, _ := makeTestUnauthenticatedRequestData()
	defer rd.Tx.Rollback()

	query := queryPropertyFromUserRow("auth_token", freshRegistrationPhoneNumber, eventWithRegistrationsID)

	var authToken sql.NullString
	row := rd.Tx.QueryRow(query)
	err := row.Scan(&authToken)
	assert.NoError(t, err)
	assert.False(t, authToken.Valid)

	rep, err := PlayerLogin(rd, &pg.PlayerLoginRequest{
		EventKey:    eventWithRegistrationsKey,
		PhoneNumber: freshRegistrationPhoneNumber,
		AuthCode:    freshRegistrationAuthCode,
	})
	assert.NoError(t, err)

	row = rd.Tx.QueryRow(query)
	err = row.Scan(&authToken)
	assert.NoError(t, err)
	assert.True(t, authToken.Valid)
	assert.Equal(t, authToken.String, rep.AuthToken)
}
