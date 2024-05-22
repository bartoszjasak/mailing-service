-- name: GetByID :one
SELECT * FROM messages
WHERE id = $1;

-- name: CreateMessage :exec
INSERT INTO messages (email, title, content, mailing_id, insert_time) VALUES ($1, $2, $3, $4, $5);

-- name: GetByMailingID :many
SELECT * FROM messages 
where mailing_id = $1;

-- name: DeleteByID :exec
DELETE FROM messages WHERE id = $1;