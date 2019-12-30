package db

import (
	"database/sql"
	"time"

	pg "github.com/escavelo/pubgolf/api/proto/pubgolf"
)

// CreatePlayer inserts a new player into the database in an unconfirmed state.
func CreatePlayer(tx *sql.Tx, eventID *string, name *string, league pg.League, phoneNumber *string,
	randCode uint32) error {
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

// CheckPlayerExistsByPhoneNumber returns true if the `eventID` + `phoneNumber` combination is valid.
func CheckPlayerExistsByPhoneNumber(tx *sql.Tx, eventID *string, phoneNumber *string) (
	bool, error) {
	playerCountRow := tx.QueryRow(`
		SELECT COUNT(*)
		FROM players
		WHERE event_id = $1
			AND phone_number = $2
		`, eventID, phoneNumber)
	var playerCount int
	err := playerCountRow.Scan(&playerCount)

	if err == sql.ErrNoRows {
		err = nil
	}

	return playerCount == 1, err
}

// GetPlayerName returns the display name for a given `playerID`.
func GetPlayerName(tx *sql.Tx, playerID *string) (string, error) {
	nameRow := tx.QueryRow(`
		SELECT name
		FROM players
		WHERE id = $2
		`, playerID)
	var name string
	err := nameRow.Scan(&name)

	if err == sql.ErrNoRows {
		err = nil
	}

	return name, err
}

// SetAuthCode updates the auth code (and expiration time) for a player's auth code.
func SetAuthCode(tx *sql.Tx, eventID *string, phoneNumber *string, authCode uint32) error {
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

// ValidateAuthCode confirms the validity (equality and expiration window) of a provided `authCode` for a given user.
func ValidateAuthCode(tx *sql.Tx, eventID *string, phoneNumber *string, authCode uint32) (bool, error) {
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
	isValid := authCodeCreatedAt != time.Time{} && time.Now().Sub(authCodeCreatedAt) < codeExpiration
	return isValid, err
}

// GenerateAuthToken unsets a user's auth code and generates a fresh auth token.
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

// GetAuthToken returns a user's auth token.
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

// ValidateAuthToken returns the corresponding event ID and player ID for a valid `authToken`. In the event of an
// invalid `authToken`, an empty string will be returned as both IDs.
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
