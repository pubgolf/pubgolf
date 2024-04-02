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

-- name: PlayerScores :many
SELECT
  v.id,
  v.name,
  COALESCE(s.value, 0)
FROM
  stages st
  JOIN venues v ON st.venue_id = v.id
  LEFT JOIN scores s ON s.stage_id = st.id
    AND s.player_id = $2
WHERE
  st.event_id = $1
ORDER BY
  st.rank ASC;

-- name: PlayerAdjustments :many
SELECT
  v.id,
  a.label,
  a.value
FROM
  stages st
  JOIN venues v ON st.venue_id = v.id
  LEFT JOIN adjustments a ON a.stage_id = st.id
    AND a.player_id = $2
WHERE
  st.event_id = $1
ORDER BY
  st.rank ASC,
  a.created_at ASC;

-- name: UpdateScore :exec
UPDATE
  scores
SET
  value = $2,
  updated_at = now()
WHERE
  id = $1;

-- name: UpdateAdjustment :exec
UPDATE
  adjustments
SET
  label = $2,
  value = $3,
  updated_at = now()
WHERE
  id = $1;

-- name: DeleteScoreForPlayerStage :exec
DELETE FROM scores
WHERE stage_id = $1
  AND player_id = $2;

-- name: DeleteAdjustmentsForPlayerStage :exec
DELETE FROM adjustments
WHERE stage_id = $1
  AND player_id = $2;

-- name: DeleteAdjustment :exec
DELETE FROM adjustments
WHERE id = $1;

