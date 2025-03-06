-- name: TrackMovieView :exec
INSERT INTO views (movie_id, user_id)
VALUES ($1, $2);

-- name: GetViewCountForMovie :one
SELECT COUNT(*) FROM views WHERE movie_id = $1;
