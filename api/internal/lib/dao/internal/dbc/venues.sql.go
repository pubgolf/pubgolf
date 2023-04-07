// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: venues.sql

package dbc

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

const venueByKey = `-- name: VenueByKey :one
SELECT
  v.id,
  v.name,
  v.address
FROM
  event_venues ev
  JOIN venues v ON ev.venue_id = v.id
WHERE
  ev.event_id = $1
  AND ev.venue_key = $2
  AND ev.deleted_at IS NULL
  AND v.deleted_at IS NULL
`

type VenueByKeyParams struct {
	EventID  models.EventID
	VenueKey models.VenueKey
}

type VenueByKeyRow struct {
	ID      models.VenueID
	Name    string
	Address string
}

func (q *Queries) VenueByKey(ctx context.Context, arg VenueByKeyParams) (VenueByKeyRow, error) {
	row := q.queryRow(ctx, q.venueByKeyStmt, venueByKey, arg.EventID, arg.VenueKey)
	var i VenueByKeyRow
	err := row.Scan(&i.ID, &i.Name, &i.Address)
	return i, err
}
