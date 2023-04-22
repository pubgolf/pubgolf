-- name: ScoringCriteriaEveryOtherVenue :many
WITH st AS (
  SELECT
    st.*,
    mod(row_number() OVER (ORDER BY st.rank ASC), 2) = 1 AS is_odd
  FROM
    stages st
  WHERE
    st.event_id = @event_id
)
SELECT
  s.player_id,
  p.name,
  count(DISTINCT (s.id)) AS num_scores,
  sum(coalesce(s.value, 0)) + sum(coalesce(a.value, 0)) AS total_points,
  sum(
    CASE WHEN a.value > 0 THEN
      a.value
    ELSE
      0
    END) AS points_from_penalties,
  sum(
    CASE WHEN a.value < 0 THEN
      a.value
    ELSE
      0
    END) AS points_from_bonuses
FROM
  scores s
  JOIN players p ON s.player_id = p.id
  JOIN st ON s.stage_id = st.id
  LEFT JOIN adjustments a ON a.stage_id = s.stage_id
    AND a.player_id = s.player_id
WHERE
  st.event_id = @event_id
  AND p.scoring_category = @scoring_category
  AND st.is_odd
GROUP BY
  s.player_id,
  p.name
ORDER BY
  num_scores DESC,
  total_points ASC,
  points_from_penalties ASC,
  points_from_bonuses DESC;

-- name: ScoringCriteriaAllVenues :many
SELECT
  s.player_id,
  p.name,
  count(DISTINCT (s.id)) AS num_scores,
  sum(coalesce(s.value, 0)) + sum(coalesce(a.value, 0)) AS total_points,
  sum(
    CASE WHEN a.value > 0 THEN
      a.value
    ELSE
      0
    END) AS points_from_penalties,
  sum(
    CASE WHEN a.value < 0 THEN
      a.value
    ELSE
      0
    END) AS points_from_bonuses
FROM
  scores s
  JOIN players p ON s.player_id = p.id
  JOIN stages st ON s.stage_id = st.id
  LEFT JOIN adjustments a ON a.stage_id = s.stage_id
    AND a.player_id = s.player_id
WHERE
  st.event_id = @event_id
  AND p.scoring_category = @scoring_category
GROUP BY
  s.player_id,
  p.name
ORDER BY
  num_scores DESC,
  total_points ASC,
  points_from_penalties ASC,
  points_from_bonuses DESC;

