package db

import (
	"database/sql"
	"time"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

func CreatePlayer(tx *sql.Tx, eventID *string, name *string, league pg.League,
	phoneNumber *string, randCode uint32) error {
	_, err := tx.Exec(`
		INSERT INTO
		players(
			event_id,
			name,
			league,
			phone_number,
			auth_code,
			auth_code_created_at,
			created_at,
			updated_at,
			phone_number_confirmed
		)
		VALUES (
			$1, $2, $3, $4, $5, NOW(), NOW(), NOW(), FALSE
		)
	`, eventID, name, league.String(), phoneNumber, randCode)
	return err
}

func CheckPlayerExists(tx *sql.Tx, eventID *string, phoneNumber *string) (
	bool, error) {
	userCountRow := tx.QueryRow(`
		SELECT COUNT(*)
		FROM players
		WHERE event_id = $1
			AND phone_number = $2
		`, eventID, phoneNumber)
	var userCount int
	err := userCountRow.Scan(&userCount)

	if err == sql.ErrNoRows {
		err = nil
	}

	return userCount == 1, err
}

func GetPlayerName(tx *sql.Tx, eventID *string, playerID *string) (
	string, error) {
	nameRow := tx.QueryRow(`
		SELECT name
		FROM players
		WHERE event_id = $1
			AND id = $2
		`, eventID, playerID)
	var name string
	err := nameRow.Scan(&name)

	if err == sql.ErrNoRows {
		err = nil
	}

	return name, err
}

func SetAuthCode(tx *sql.Tx, eventID *string, phoneNumber *string,
	authCode uint32) error {
	_, err := tx.Exec(`
		UPDATE players
		SET
			auth_code = $3,
			auth_code_created_at = NOW(),
			auth_token = NULL
		WHERE event_id = $1
			AND phone_number = $2
		`, eventID, phoneNumber, authCode)
	return err
}

func ValidateAuthCode(tx *sql.Tx, eventID *string, phoneNumber *string,
	authCode uint32) (bool, error) {
	authCodeCreatedAtRow := tx.QueryRow(`
		SELECT auth_code_created_at
		FROM players
		WHERE event_id = $1
			AND phone_number = $2
			AND auth_code = $3
		`, eventID, phoneNumber, authCode)
	var authCodeCreatedAt time.Time
	err := authCodeCreatedAtRow.Scan(&authCodeCreatedAt)

	if err == sql.ErrNoRows {
		err = nil
	}

	if err != nil {
		return false, err
	}
	codeExpiration, err := time.ParseDuration("10m")
	isValid := authCodeCreatedAt != time.Time{} &&
		time.Now().Sub(authCodeCreatedAt) < codeExpiration
	return isValid, err
}

func GenerateAuthToken(tx *sql.Tx, eventID *string, phoneNumber *string) error {
	_, err := tx.Exec(`
		UPDATE players
		SET
			phone_number_confirmed = true,
			auth_token = uuid_generate_v4(),
			auth_code = NULL,
			auth_code_created_at = NULL,
			updated_at = NOW()
		WHERE event_id = $1 AND phone_number = $2
	`, eventID, phoneNumber)
	return err
}

func GetAuthToken(tx *sql.Tx, eventID *string, phoneNumber *string) (string, error) {
	authTokenRow := tx.QueryRow(`
		SELECT auth_token
		FROM players
		WHERE event_id = $1
			AND phone_number = $2
		`, eventID, phoneNumber)
	var authToken string
	err := authTokenRow.Scan(&authToken)
	return authToken, err
}

func ValidateAuthToken(tx *sql.Tx, authToken *string) (string, string, error) {
	var eventID, playerID string
	row := tx.QueryRow(`
		SELECT event_id, id
		FROM players
		WHERE auth_token = $1
		`, authToken)
	err := row.Scan(&eventID, &playerID)
	if err == sql.ErrNoRows {
		err = nil
	}
	return eventID, playerID, err
}
