-- name: CreateUser :one
INSERT into users(
    google_id,
    name,
    email,
    profile_image
) VALUES (
 $1, $2, $3, $4
) RETURNING id;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $1
WHERE email = $2
RETURNING name;