-- name: CreateChirp :one
INSERT INTO chirps(body, user_id)
VALUES($1, $2)
RETURNING *;

-- name: GetChirpByID :one
SELECT * FROM chirps Where id = $1;

-- name: GetChirpByUserID :many
SELECT * FROM chirps Where user_id = $1 ORDER BY created_at ASC;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: DeleteChirpByUserID :exec
DELETE FROM chirps WHERE id = $1 AND user_id = $2;

