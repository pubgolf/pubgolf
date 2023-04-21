-- name: ScoreByPlayerStage :one
SELECT
  id,
  value
FROM
  scores
WHERE
  stage_id = $1
  AND player_id = $2
  AND deleted_at IS NULL;

-- name: AdjustmentsByPlayerStage :many
SELECT
  id,
  label,
  value
FROM
  adjustments
WHERE
  stage_id = $1
  AND player_id = $2
  AND deleted_at IS NULL;

-- name: CreateScore :exec
INSERT INTO scores(stage_id, player_id, value, updated_at)
  VALUES ($1, $2, $3, now());

-- name: CreateAdjustment :exec
INSERT INTO adjustments(stage_id, player_id, label, value)
  VALUES ($1, $2, $3, $4);

