-- name: RuleItemsByStageIDs :many
SELECT
  id,
  stages_id,
  content,
  item_type,
  audiences,
  rank
FROM
  rule_items
WHERE
  stages_id = ANY(@stage_ids::uuid[])
  AND deleted_at IS NULL
ORDER BY
  stages_id,
  rank ASC;

-- name: DeleteRuleItemsByStageID :exec
UPDATE
  rule_items
SET
  deleted_at = now()
WHERE
  stages_id = $1
  AND deleted_at IS NULL;

-- name: CreateRuleItem :exec
INSERT INTO rule_items (stages_id, content, item_type, audiences, rank)
  VALUES ($1, $2, $3, $4, $5);
