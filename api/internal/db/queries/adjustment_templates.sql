-- name: CreateAdjustmentTemplate :one
INSERT INTO adjustment_templates(label, value, rank, stage_id, event_id, deleted_at)
  VALUES (@label, @value, @rank, @stage_id, @event_id, @deleted_at)
RETURNING
  id;

-- name: AdjustmentTemplatesByEventID :many
SELECT
  at.id,
  at.value,
  at.label
FROM
  adjustment_templates at
WHERE
  at.deleted_at IS NULL
  AND at.event_id = @event_id
ORDER BY
  at.rank;

-- name: AdjustmentTemplatesByStageID :many
SELECT
  at.id,
  at.value,
  at.label
FROM
  adjustment_templates at
WHERE
  at.deleted_at IS NULL
  AND at.stage_id = @stage_id
ORDER BY
  at.rank;

-- name: EventAdjustmentTemplates :many
SELECT
  at.id,
  at.value,
  at.label,
  at.rank,
  at.stage_id,
  CASE WHEN at.deleted_at IS NULL THEN
    TRUE
  ELSE
    FALSE
  END AS is_visible
FROM
  adjustment_templates at
  LEFT JOIN stages s ON at.stage_id = s.id
WHERE
  -- Don't check `at.deleted_at IS NULL` to include soft-deleted templates.
((at.event_id = @event_id)
    OR (s.deleted_at IS NULL
      AND s.event_id = @event_id))
ORDER BY
  s.rank ASC,
  at.rank ASC;

-- name: UpdateAdjustmentTemplate :exec
UPDATE
  adjustment_templates
SET
  label = @label,
  value = @value,
  rank = @rank,
  stage_id = @stage_id,
  event_id = @event_id,
  deleted_at = @deleted_at
WHERE
  id = @id;

