-- name: CreateArtist :one
INSERT INTO artists (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE 
SET name = EXCLUDED.name
RETURNING *;

-- name: GetArtistByName :one
SELECT * FROM artists WHERE name = $1;

-- name: GetArtistByID :one
SELECT * FROM artists WHERE id = $1;

-- name: UpsertArtist :one
INSERT INTO artists (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;
