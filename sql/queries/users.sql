-- name: CreateUser :one
INSERT INTO users(id, email, password)
VALUES($1, $2, $3)
RETURNING *;


-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserByEmailAndPassword :exec
UPDATE users
SET email = $1, password = $2
WHERE id = $3;

-- name: UpdateChirpyIsRed :exec
UPDATE users
SET is_chirpy_red = TRUE, updated_at = $1
WHERE id = $2;

-- name: DeleteAllUsers :exec
DELETE FROM users;
