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

-- name: EventScores :many
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
  s.created_at ASC;

-- name: EventAdjustments :many
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
  a.created_at ASC;

