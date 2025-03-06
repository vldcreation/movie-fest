-- name: GetMovies :many
SELECT m.*, 
       array_agg(DISTINCT g.name) AS genres,
       array_agg(DISTINCT a.name) AS artists,
       COUNT(DISTINCT v.id) AS view_count,
       COUNT(DISTINCT vt.user_id) AS vote_count
FROM movies m
LEFT JOIN movie_genres mg ON m.id = mg.movie_id
LEFT JOIN genres g ON mg.genre_id = g.id
LEFT JOIN movie_artists ma ON m.id = ma.movie_id
LEFT JOIN artists a ON ma.artist_id = a.id
LEFT JOIN views v ON m.id = v.movie_id
LEFT JOIN votes vt ON m.id = vt.movie_id
WHERE
    ($1::text = '' OR m.title ILIKE '%' || $1 || '%') AND
    ($2::text = '' OR m.description ILIKE '%' || $2 || '%') AND
    ($3::text = '' OR g.name ILIKE '%' || $3 || '%') AND
    ($4::text = '' OR a.name ILIKE '%' || $4 || '%')
GROUP BY m.id
ORDER BY m.created_at DESC
LIMIT $5 OFFSET $6;

-- name: GetMovieByID :one
SELECT m.*, 
       array_agg(DISTINCT g.name) AS genres,
       array_agg(DISTINCT a.name) AS artists,
       COUNT(DISTINCT v.id) AS view_count,
       COUNT(DISTINCT vt.user_id) AS vote_count
FROM movies m
LEFT JOIN movie_genres mg ON m.id = mg.movie_id
LEFT JOIN genres g ON mg.genre_id = g.id
LEFT JOIN movie_artists ma ON m.id = ma.movie_id
LEFT JOIN artists a ON ma.artist_id = a.id
LEFT JOIN views v ON m.id = v.movie_id
LEFT JOIN votes vt ON m.id = vt.movie_id
WHERE m.id = $1
GROUP BY m.id;

-- name: CreateMovie :one
INSERT INTO movies (title, description, duration, watch_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2, description = $3, duration = $4, watch_url = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = $1;

-- name: AddGenreToMovie :exec
INSERT INTO movie_genres (movie_id, genre_id)
VALUES ($1, $2);

-- name: AddArtistToMovie :exec
INSERT INTO movie_artists (movie_id, artist_id)
VALUES ($1, $2);

-- name: GetMostViewedMovies :many
SELECT m.*, COUNT(v.id) AS view_count
FROM movies m
LEFT JOIN views v ON m.id = v.movie_id
GROUP BY m.id
ORDER BY view_count DESC
LIMIT $1 OFFSET $2;

-- name: GetMostViewedGenres :many
SELECT g.name, COUNT(v.id) AS view_count
FROM genres g
JOIN movie_genres mg ON g.id = mg.genre_id
JOIN views v ON mg.movie_id = v.movie_id
GROUP BY g.id
ORDER BY view_count DESC
LIMIT $1 OFFSET $2;

-- name: RemoveAllGenresFromMovie :exec
DELETE FROM movie_genres WHERE movie_id = $1;

-- name: RemoveAllArtistsFromMovie :exec
DELETE FROM movie_artists WHERE movie_id = $1;
