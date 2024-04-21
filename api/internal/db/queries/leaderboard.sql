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
      s.player_id,
      p.name,
      count(DISTINCT (s.id)) AS num_scores,
      sum(s.value) AS total_points,
      0 AS points_from_penalties,
      0 AS points_from_bonuses
    FROM
      scores s
      JOIN players p ON s.player_id = p.id
      JOIN event_players ep ON p.id = ep.player_id
      JOIN st ON s.stage_id = st.id
    WHERE
      s.deleted_at IS NULL
      AND s.is_verified IS TRUE
      AND p.deleted_at IS NULL
      AND ep.deleted_at IS NULL
      AND ep.event_id = @event_id
      AND ep.scoring_category = @scoring_category
      AND st.deleted_at IS NULL
      AND st.event_id = @event_id
      AND st.is_odd
    GROUP BY
      s.player_id,
      p.name)
  UNION (
    -- Get adjustment contributions.
    SELECT
      s.player_id,
      p.name,
      count(DISTINCT (s.id)) AS num_scores,
      sum(a.value) AS total_points,
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
      JOIN event_players ep ON p.id = ep.player_id
      JOIN st ON s.stage_id = st.id
      LEFT JOIN adjustments a ON a.stage_id = s.stage_id
        AND a.player_id = s.player_id
    WHERE
      s.deleted_at IS NULL
      AND s.is_verified IS TRUE
      AND p.deleted_at IS NULL
      AND ep.deleted_at IS NULL
      AND ep.event_id = @event_id
      AND ep.scoring_category = @scoring_category
      AND st.deleted_at IS NULL
      AND st.event_id = @event_id
      AND st.is_odd
      AND a.deleted_at IS NULL
    GROUP BY
      s.player_id,
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
      s.player_id,
      p.name,
      count(DISTINCT (s.id)) AS num_scores,
      sum(s.value) AS total_points,
      0 AS points_from_penalties,
      0 AS points_from_bonuses
    FROM
      scores s
      JOIN players p ON s.player_id = p.id
      JOIN event_players ep ON p.id = ep.player_id
      JOIN stages st ON s.stage_id = st.id
    WHERE
      s.deleted_at IS NULL
      AND s.is_verified IS TRUE
      AND p.deleted_at IS NULL
      AND ep.deleted_at IS NULL
      AND ep.event_id = @event_id
      AND ep.scoring_category = @scoring_category
      AND st.deleted_at IS NULL
      AND st.event_id = @event_id
    GROUP BY
      s.player_id,
      p.name)
  UNION (
    -- Get adjustment contributions.
    SELECT
      s.player_id,
      p.name,
      count(DISTINCT (s.id)) AS num_scores,
      sum(a.value) AS total_points,
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
      JOIN event_players ep ON p.id = ep.player_id
      JOIN stages st ON s.stage_id = st.id
      LEFT JOIN adjustments a ON a.stage_id = s.stage_id
        AND a.player_id = s.player_id
    WHERE
      s.deleted_at IS NULL
      AND s.is_verified IS TRUE
      AND p.deleted_at IS NULL
      AND ep.deleted_at IS NULL
      AND ep.event_id = @event_id
      AND ep.scoring_category = @scoring_category
      AND st.deleted_at IS NULL
      AND st.event_id = @event_id
      AND a.deleted_at IS NULL
    GROUP BY
      s.player_id,
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

