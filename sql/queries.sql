-- name: CreateFile :exec
INSERT INTO files (id, original_name, saved_name, file_size) 
VALUES (?, ?, ?, ?);

-- name: ListFiles :many
SELECT id, original_name, saved_name, file_size, created_at 
FROM files 
ORDER BY created_at DESC;

-- name: GetFile :one
SELECT id, original_name, saved_name, file_size, created_at 
FROM files 
WHERE id = ? LIMIT 1;

-- name: GetOldFiles :many
SELECT id, saved_name 
FROM files 
WHERE created_at < datetime('now', '-24 hours');

-- name: DeleteFile :exec
DELETE FROM files WHERE id = ?;

-- name: GetFileByID :one
SELECT id, original_name, saved_name, file_size, created_at
FROM files
WHERE id = ? LIMIT 1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = ?;