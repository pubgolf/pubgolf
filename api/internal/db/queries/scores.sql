-- name: UpsertScore :exec
INSERT INTO scores(stage_id, player_id, value, is_verified)
  VALUES (@stage_id, @player_id, @value, @is_verified)
ON CONFLICT (stage_id, player_id)
  DO UPDATE SET
    value = EXCLUDED.value, is_verified = EXCLUDED.is_verified, updated_at = now(), deleted_at = NULL;

-- name: CreateAdjustment :exec
INSERT INTO adjustments(stage_id, player_id, label, value)
  VALUES ($1, $2, $3, $4);

-- name: CreateAdjustmentWithTemplate :exec
INSERT INTO adjustments(stage_id, player_id, label, value, adjustment_template_id)
  VALUES ($1, $2, $3, $4, $5);

-- name: ScoreByPlayerStage :one
SELECT
  id,
  value,
  is_verified
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
  value,
  adjustment_template_id
FROM
  adjustments
WHERE
  stage_id = $1
  AND player_id = $2
  AND deleted_at IS NULL;

-- name: EventScores :many
SELECT
  s.stage_id,
  s.player_id,
  s.id AS score_id,
  s.value,
  s.is_verified
FROM
  scores s
  JOIN stages st ON s.stage_id = st.id
WHERE
  s.deleted_at IS NULL
  AND (s.is_verified = FALSE
    OR s.is_verified = @include_verified)
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
  COALESCE(s.value, 0),
  COALESCE(s.is_verified, TRUE)
FROM
  stages st
  JOIN venues v ON st.venue_id = v.id
  LEFT JOIN scores s ON s.stage_id = st.id
    AND s.player_id = @player_id
WHERE
  st.deleted_at IS NULL
  AND st.event_id = @event_id
  AND v.deleted_at IS NULL
  AND s.deleted_at IS NULL
  AND (s.is_verified = TRUE
    OR s.is_verified != @include_unverified)
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
  st.deleted_at IS NULL
  AND st.event_id = $1
  AND v.deleted_at IS NULL
  AND a.deleted_at IS NULL
ORDER BY
  st.rank ASC,
  a.created_at ASC;

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

