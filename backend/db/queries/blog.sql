-- name: CreateBlog :one
INSERT INTO blogs(
    title,
    content,
    author
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetBlogByID :one
SELECT 
b.id, 
b.title, 
b.content,
b.created_at, 
b.version, 
u.name, 
u.id AS author_id
FROM blogs b
INNER JOIN users u
ON u.id = b.author
WHERE b.id = $1;

-- name: GetBlogs :many
SELECT 
b.id, 
b.title, 
b.content,
b.created_at, 
b.version, 
u.name, 
u.id AS author_id
FROM blogs b
INNER JOIN users u
ON u.id = b.author
ORDER BY created_at DESC, b.id ASC
LIMIT $1 OFFSET $2;

-- name: GetBlogCount :one
SELECT COUNT(*) FROM blogs;

-- name: UpdateBlog :one
UPDATE blogs
SET title = $1,
    content = $2,
    version = version + 1
WHERE id = $3
RETURNING *;

-- name: DeleteBlog :exec
DELETE FROM blogs
WHERE id = $1 AND author = $2;