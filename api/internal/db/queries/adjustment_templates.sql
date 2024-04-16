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

