-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByID :one
SELECT * FROM chirps
WHERE id = $1; 

-- name: GetChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;