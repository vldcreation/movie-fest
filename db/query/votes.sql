-- name: VoteMovie :exec
INSERT INTO votes (user_id, movie_id)
VALUES ($1, $2)
ON CONFLICT (user_id, movie_id) DO NOTHING;

-- name: UnvoteMovie :exec
DELETE FROM votes
WHERE user_id = $1 AND movie_id = $2;

-- name: GetUserVotedMovies :many
SELECT m.*
FROM movies m
JOIN votes v ON m.id = v.movie_id
WHERE v.user_id = $1;

-- name: GetMostVotedMovies :many
SELECT m.*, COUNT(v.user_id) AS vote_count
FROM movies m
LEFT JOIN votes v ON m.id = v.movie_id
GROUP BY m.id
ORDER BY vote_count DESC
LIMIT $1;
