package db

import "database/sql"

func GetEventID(tx *sql.Tx, eventKey *string) (string, error) {
	eventIDRow := tx.QueryRow(`
		SELECT id
		FROM events
		WHERE key = $1
		`, eventKey)
	var eventID string
	err := eventIDRow.Scan(&eventID)
	return eventID, err
}
