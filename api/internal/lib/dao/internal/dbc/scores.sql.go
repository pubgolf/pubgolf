// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: scores.sql

package dbc

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

const adjustmentsByPlayerStage = `-- name: AdjustmentsByPlayerStage :many
SELECT
  id,
  label,
  value
FROM
  adjustments
WHERE
  stage_id = $1
  AND player_id = $2
  AND deleted_at IS NULL
`

type AdjustmentsByPlayerStageParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
}

type AdjustmentsByPlayerStageRow struct {
	ID    models.AdjustmentID
	Label string
	Value int32
}

func (q *Queries) AdjustmentsByPlayerStage(ctx context.Context, arg AdjustmentsByPlayerStageParams) ([]AdjustmentsByPlayerStageRow, error) {
	rows, err := q.query(ctx, q.adjustmentsByPlayerStageStmt, adjustmentsByPlayerStage, arg.StageID, arg.PlayerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdjustmentsByPlayerStageRow
	for rows.Next() {
		var i AdjustmentsByPlayerStageRow
		if err := rows.Scan(&i.ID, &i.Label, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createAdjustment = `-- name: CreateAdjustment :exec
INSERT INTO adjustments(stage_id, player_id, label, value)
  VALUES ($1, $2, $3, $4)
`

type CreateAdjustmentParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
	Label    string
	Value    int32
}

func (q *Queries) CreateAdjustment(ctx context.Context, arg CreateAdjustmentParams) error {
	_, err := q.exec(ctx, q.createAdjustmentStmt, createAdjustment,
		arg.StageID,
		arg.PlayerID,
		arg.Label,
		arg.Value,
	)
	return err
}

const createScore = `-- name: CreateScore :exec
INSERT INTO scores(stage_id, player_id, value, updated_at)
  VALUES ($1, $2, $3, now())
`

type CreateScoreParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
	Value    uint32
}

func (q *Queries) CreateScore(ctx context.Context, arg CreateScoreParams) error {
	_, err := q.exec(ctx, q.createScoreStmt, createScore, arg.StageID, arg.PlayerID, arg.Value)
	return err
}

const deleteAdjustment = `-- name: DeleteAdjustment :exec
DELETE FROM adjustments
WHERE id = $1
`

func (q *Queries) DeleteAdjustment(ctx context.Context, id models.AdjustmentID) error {
	_, err := q.exec(ctx, q.deleteAdjustmentStmt, deleteAdjustment, id)
	return err
}

const deleteAdjustmentsForPlayerStage = `-- name: DeleteAdjustmentsForPlayerStage :exec
DELETE FROM adjustments
WHERE stage_id = $1
  AND player_id = $2
`

type DeleteAdjustmentsForPlayerStageParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
}

func (q *Queries) DeleteAdjustmentsForPlayerStage(ctx context.Context, arg DeleteAdjustmentsForPlayerStageParams) error {
	_, err := q.exec(ctx, q.deleteAdjustmentsForPlayerStageStmt, deleteAdjustmentsForPlayerStage, arg.StageID, arg.PlayerID)
	return err
}

const deleteScoreForPlayerStage = `-- name: DeleteScoreForPlayerStage :exec
DELETE FROM scores
WHERE stage_id = $1
  AND player_id = $2
`

type DeleteScoreForPlayerStageParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
}

func (q *Queries) DeleteScoreForPlayerStage(ctx context.Context, arg DeleteScoreForPlayerStageParams) error {
	_, err := q.exec(ctx, q.deleteScoreForPlayerStageStmt, deleteScoreForPlayerStage, arg.StageID, arg.PlayerID)
	return err
}

const eventAdjustments = `-- name: EventAdjustments :many
SELECT
  s.stage_id,
  s.player_id,
  a.id AS adjustment_id,
  a.label,
  a.value
FROM
  scores s
  JOIN stages st ON s.stage_id = st.id
  JOIN adjustments a ON a.stage_id = s.stage_id
    AND a.player_id = s.player_id
WHERE
  s.deleted_at IS NULL
  AND st.deleted_at IS NULL
  AND a.deleted_at IS NULL
  AND st.event_id = $1
ORDER BY
  st.rank ASC,
  s.created_at ASC,
  a.created_at ASC
`

type EventAdjustmentsRow struct {
	StageID      models.StageID
	PlayerID     models.PlayerID
	AdjustmentID models.DatabaseULID
	Label        string
	Value        int32
}

func (q *Queries) EventAdjustments(ctx context.Context, eventID models.EventID) ([]EventAdjustmentsRow, error) {
	rows, err := q.query(ctx, q.eventAdjustmentsStmt, eventAdjustments, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventAdjustmentsRow
	for rows.Next() {
		var i EventAdjustmentsRow
		if err := rows.Scan(
			&i.StageID,
			&i.PlayerID,
			&i.AdjustmentID,
			&i.Label,
			&i.Value,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const eventScores = `-- name: EventScores :many
SELECT
  s.stage_id,
  s.player_id,
  s.id AS score_id,
  s.value
FROM
  scores s
  JOIN stages st ON s.stage_id = st.id
WHERE
  s.deleted_at IS NULL
  AND st.deleted_at IS NULL
  AND st.event_id = $1
ORDER BY
  st.rank ASC,
  s.created_at ASC
`

type EventScoresRow struct {
	StageID  models.StageID
	PlayerID models.PlayerID
	ScoreID  models.DatabaseULID
	Value    uint32
}

func (q *Queries) EventScores(ctx context.Context, eventID models.EventID) ([]EventScoresRow, error) {
	rows, err := q.query(ctx, q.eventScoresStmt, eventScores, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EventScoresRow
	for rows.Next() {
		var i EventScoresRow
		if err := rows.Scan(
			&i.StageID,
			&i.PlayerID,
			&i.ScoreID,
			&i.Value,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const scoreByPlayerStage = `-- name: ScoreByPlayerStage :one
SELECT
  id,
  value
FROM
  scores
WHERE
  stage_id = $1
  AND player_id = $2
  AND deleted_at IS NULL
`

type ScoreByPlayerStageParams struct {
	StageID  models.StageID
	PlayerID models.PlayerID
}

type ScoreByPlayerStageRow struct {
	ID    models.ScoreID
	Value uint32
}

func (q *Queries) ScoreByPlayerStage(ctx context.Context, arg ScoreByPlayerStageParams) (ScoreByPlayerStageRow, error) {
	row := q.queryRow(ctx, q.scoreByPlayerStageStmt, scoreByPlayerStage, arg.StageID, arg.PlayerID)
	var i ScoreByPlayerStageRow
	err := row.Scan(&i.ID, &i.Value)
	return i, err
}

const updateAdjustment = `-- name: UpdateAdjustment :exec
UPDATE
  adjustments
SET
  label = $2,
  value = $3,
  updated_at = now()
WHERE
  id = $1
`

type UpdateAdjustmentParams struct {
	ID    models.AdjustmentID
	Label string
	Value int32
}

func (q *Queries) UpdateAdjustment(ctx context.Context, arg UpdateAdjustmentParams) error {
	_, err := q.exec(ctx, q.updateAdjustmentStmt, updateAdjustment, arg.ID, arg.Label, arg.Value)
	return err
}

const updateScore = `-- name: UpdateScore :exec
UPDATE
  scores
SET
  value = $2,
  updated_at = now()
WHERE
  id = $1
`

type UpdateScoreParams struct {
	ID    models.ScoreID
	Value uint32
}

func (q *Queries) UpdateScore(ctx context.Context, arg UpdateScoreParams) error {
	_, err := q.exec(ctx, q.updateScoreStmt, updateScore, arg.ID, arg.Value)
	return err
}
