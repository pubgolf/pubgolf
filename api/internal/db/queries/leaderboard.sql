-- name: ScoringCriteriaEveryOtherVenue :many
WITH st AS (
  -- Replaces the stages table in the later section with only the odd numbered stages.
  SELECT
    st.*,
    mod(row_number() OVER (ORDER BY st.rank ASC), 2) = 1 AS is_odd
  FROM
    stages st
  WHERE
    st.deleted_at IS NULL
    AND st.event_id = @event_id
),
separated AS ((
    -- Get score contributions.
    SELECT
      p.id AS player_id,
      p.name,
      coalesce(count(DISTINCT (s.id)), 0)::bigint AS num_scores,
      coalesce(sum(s.value), 0)::bigint AS total_points,
      0 AS points_from_penalties,
      0 AS points_from_bonuses
    FROM
      players p
    LEFT JOIN event_players ep ON p.id = ep.player_id
    LEFT JOIN scores s ON s.player_id = p.id
      AND s.is_verified IS TRUE
    LEFT JOIN st ON s.stage_id = st.id
      AND st.event_id = ep.event_id
  WHERE
    p.deleted_at IS NULL
    AND ep.deleted_at IS NULL
    AND ep.event_id = @event_id
    AND ep.scoring_category = @scoring_category
    AND s.deleted_at IS NULL
    AND st.deleted_at IS NULL
    AND st.is_odd
  GROUP BY
    p.id,
    p.name)
UNION (
  -- Get adjustment contributions.
  SELECT
    p.id AS player_id,
    p.name,
    coalesce(count(DISTINCT (s.id)), 0)::bigint AS num_scores,
    coalesce(sum(a.value), 0)::bigint AS total_points,
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
    players p
    LEFT JOIN event_players ep ON p.id = ep.player_id
    LEFT JOIN scores s ON s.player_id = p.id
      AND s.is_verified IS TRUE
    LEFT JOIN st ON s.stage_id = st.id
      AND st.event_id = ep.event_id
    LEFT JOIN adjustments a ON a.stage_id = s.stage_id
      AND a.player_id = s.player_id
  WHERE
    p.deleted_at IS NULL
    AND ep.deleted_at IS NULL
    AND ep.event_id = @event_id
    AND ep.scoring_category = @scoring_category
    AND s.deleted_at IS NULL
    AND st.deleted_at IS NULL
    AND st.is_odd
    AND a.deleted_at IS NULL
  GROUP BY
    p.id,
    p.name))
SELECT
  player_id,
  name,
  num_scores,
  SUM(total_points) AS total_points,
  SUM(points_from_penalties) AS points_from_penalties,
  SUM(points_from_bonuses) AS points_from_bonuses
FROM
  separated
GROUP BY
  player_id,
  name,
  num_scores
ORDER BY
  num_scores DESC,
  total_points ASC,
  points_from_penalties ASC,
  points_from_bonuses DESC;

-- name: ScoringCriteriaAllVenues :many
WITH separated AS ((
    -- Get score contributions.
    SELECT
      p.id AS player_id,
      p.name,
      coalesce(count(DISTINCT (s.id)), 0)::bigint AS num_scores,
      coalesce(sum(s.value), 0)::bigint AS total_points,
      0 AS points_from_penalties,
      0 AS points_from_bonuses
    FROM
      players p
    LEFT JOIN event_players ep ON p.id = ep.player_id
    LEFT JOIN scores s ON s.player_id = p.id
      AND s.is_verified IS TRUE
    LEFT JOIN stages st ON s.stage_id = st.id
      AND st.event_id = ep.event_id
  WHERE
    p.deleted_at IS NULL
    AND ep.deleted_at IS NULL
    AND ep.event_id = @event_id
    AND ep.scoring_category = @scoring_category
    AND s.deleted_at IS NULL
    AND st.deleted_at IS NULL
  GROUP BY
    p.id,
    p.name)
UNION (
  -- Get adjustment contributions.
  SELECT
    p.id AS player_id,
    p.name,
    coalesce(count(DISTINCT (s.id)), 0)::bigint AS num_scores,
    coalesce(sum(a.value), 0)::bigint AS total_points,
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
    players p
    LEFT JOIN event_players ep ON p.id = ep.player_id
    LEFT JOIN scores s ON s.player_id = p.id
      AND s.is_verified IS TRUE
    LEFT JOIN stages st ON s.stage_id = st.id
      AND st.event_id = ep.event_id
    LEFT JOIN adjustments a ON a.stage_id = s.stage_id
      AND a.player_id = s.player_id
  WHERE
    p.deleted_at IS NULL
    AND ep.deleted_at IS NULL
    AND ep.event_id = @event_id
    AND ep.scoring_category = @scoring_category
    AND s.deleted_at IS NULL
    AND st.deleted_at IS NULL
    AND a.deleted_at IS NULL
  GROUP BY
    p.id,
    p.name))
SELECT
  player_id,
  name,
  num_scores,
  SUM(total_points) AS total_points,
  SUM(points_from_penalties) AS points_from_penalties,
  SUM(points_from_bonuses) AS points_from_bonuses
FROM
  separated
GROUP BY
  player_id,
  name,
  num_scores
ORDER BY
  num_scores DESC,
  total_points ASC,
  points_from_penalties ASC,
  points_from_bonuses DESC;

-- name: UnverifiedScoreCountEveryOtherVenue :many
WITH st AS (
  -- Replaces the stages table in the later section with only the odd numbered stages.
  SELECT
    st.*,
    mod(row_number() OVER (ORDER BY st.rank ASC), 2) = 1 AS is_odd
  FROM
    stages st
  WHERE
    st.deleted_at IS NULL
    AND st.event_id = @event_id
)
SELECT
  p.id AS player_id,
  SUM(
    CASE WHEN s.is_verified IS FALSE THEN
      1
    ELSE
      0
    END) AS count
FROM
  players p
  JOIN event_players ep ON p.id = ep.player_id
  JOIN st ON st.event_id = ep.event_id
  LEFT JOIN scores s ON s.player_id = p.id
    AND s.stage_id = st.id
WHERE
  p.deleted_at IS NULL
  AND ep.deleted_at IS NULL
  AND ep.event_id = @event_id
  AND ep.scoring_category = @scoring_category
  AND s.deleted_at IS NULL
  AND st.deleted_at IS NULL
  AND st.is_odd
GROUP BY
  p.id;

-- name: UnverifiedScoreCountAllVenues :many
SELECT
  p.id AS player_id,
  SUM(
    CASE WHEN s.is_verified IS FALSE THEN
      1
    ELSE
      0
    END) AS count
FROM
  players p
  JOIN event_players ep ON p.id = ep.player_id
  JOIN stages st ON st.event_id = ep.event_id
  LEFT JOIN scores s ON s.player_id = p.id
    AND s.stage_id = st.id
WHERE
  p.deleted_at IS NULL
  AND ep.deleted_at IS NULL
  AND ep.event_id = @event_id
  AND ep.scoring_category = @scoring_category
  AND s.deleted_at IS NULL
  AND st.deleted_at IS NULL
GROUP BY
  p.id;

