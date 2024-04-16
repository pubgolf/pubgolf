// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: adjustment_templates.sql

package dbc

import (
	"context"

	"github.com/pubgolf/pubgolf/api/internal/lib/models"
)

const adjustmentTemplatesByEventID = `-- name: AdjustmentTemplatesByEventID :many
SELECT
  at.id,
  at.value,
  at.label
FROM
  adjustment_templates at
WHERE
  at.deleted_at IS NULL
  AND at.event_id = $1
ORDER BY
  at.rank
`

type AdjustmentTemplatesByEventIDRow struct {
	ID    models.AdjustmentTemplateID
	Value int32
	Label string
}

func (q *Queries) AdjustmentTemplatesByEventID(ctx context.Context, eventID models.EventID) ([]AdjustmentTemplatesByEventIDRow, error) {
	rows, err := q.query(ctx, q.adjustmentTemplatesByEventIDStmt, adjustmentTemplatesByEventID, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdjustmentTemplatesByEventIDRow
	for rows.Next() {
		var i AdjustmentTemplatesByEventIDRow
		if err := rows.Scan(&i.ID, &i.Value, &i.Label); err != nil {
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

const adjustmentTemplatesByStageID = `-- name: AdjustmentTemplatesByStageID :many
SELECT
  at.id,
  at.value,
  at.label
FROM
  adjustment_templates at
WHERE
  at.deleted_at IS NULL
  AND at.stage_id = $1
ORDER BY
  at.rank
`

type AdjustmentTemplatesByStageIDRow struct {
	ID    models.AdjustmentTemplateID
	Value int32
	Label string
}

func (q *Queries) AdjustmentTemplatesByStageID(ctx context.Context, stageID models.StageID) ([]AdjustmentTemplatesByStageIDRow, error) {
	rows, err := q.query(ctx, q.adjustmentTemplatesByStageIDStmt, adjustmentTemplatesByStageID, stageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AdjustmentTemplatesByStageIDRow
	for rows.Next() {
		var i AdjustmentTemplatesByStageIDRow
		if err := rows.Scan(&i.ID, &i.Value, &i.Label); err != nil {
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
