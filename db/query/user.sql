-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2);

-- name: GetUserRoles :one
SELECT r.id, r.name
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;
