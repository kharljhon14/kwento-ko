-- name: AddBlogTags :exec
INSERT INTO blog_tags (blog_id, tag_id)
SELECT $1, unnest($2::uuid[]);

-- name: RemoveBlogTags :exec
DELETE FROM blog_tags
WHERE blog_id = $1
    AND tag_id = ANY($2::uuid[]);


-- name: GetBlogTags :many
SELECT  t.name
FROM blog_tags b
INNER JOIN tags t 
ON t.id = b.tag_id
WHERE b.blog_id = $1;
