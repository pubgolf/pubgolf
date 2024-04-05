// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: auth.sql

package dbc

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

const generateAuthToken = `-- name: GenerateAuthToken :one
INSERT INTO auth_tokens(player_id)
SELECT
  id
FROM
  players
WHERE
  phone_number = $1
RETURNING
  player_id,
  id AS auth_token
`

type GenerateAuthTokenRow struct {
	PlayerID  models.PlayerID
	AuthToken models.AuthToken
}

func (q *Queries) GenerateAuthToken(ctx context.Context, phoneNumber models.PhoneNum) (GenerateAuthTokenRow, error) {
	row := q.queryRow(ctx, q.generateAuthTokenStmt, generateAuthToken, phoneNumber)
	var i GenerateAuthTokenRow
	err := row.Scan(&i.PlayerID, &i.AuthToken)
	return i, err
}