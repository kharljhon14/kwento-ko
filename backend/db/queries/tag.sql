-- name: CreateTag :one
INSERT INTO tags(
    name
) VALUES( 
    $1
) RETURNING *;

-- name: UpdateTag :one
UPDATE tags
SET name = $1
WHERE id = $2
RETURNING *;

-- name: GetTag :one
SELECT * FROM tags
WHERE id = $1;

-- name: GetTags :many
SELECT * FROM tags
LIMIT $1 OFFSET $2;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1;
