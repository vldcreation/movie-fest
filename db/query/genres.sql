-- name: CreateGenre :one
INSERT INTO genres (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE 
SET name = EXCLUDED.name
RETURNING *;

-- name: GetGenreByName :one
SELECT * FROM genres WHERE name = $1;

-- name: GetGenreByID :one
SELECT * FROM genres WHERE id = $1;

-- name: UpsertGenre :one
INSERT INTO genres (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;