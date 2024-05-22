// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc

import (
	"context"
	"time"
)

const createMessage = `-- name: CreateMessage :exec
INSERT INTO messages (email, title, content, mailing_id, insert_time) VALUES ($1, $2, $3, $4, $5)
`

type CreateMessageParams struct {
	Email      string
	Title      string
	Content    string
	MailingID  int32
	InsertTime time.Time
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) error {
	_, err := q.db.ExecContext(ctx, createMessage,
		arg.Email,
		arg.Title,
		arg.Content,
		arg.MailingID,
		arg.InsertTime,
	)
	return err
}

const deleteByID = `-- name: DeleteByID :exec
DELETE FROM messages WHERE id = $1
`

func (q *Queries) DeleteByID(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteByID, id)
	return err
}

const getByID = `-- name: GetByID :one
SELECT id, email, title, content, mailing_id, insert_time FROM messages
WHERE id = $1
`

func (q *Queries) GetByID(ctx context.Context, id int32) (Message, error) {
	row := q.db.QueryRowContext(ctx, getByID, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Title,
		&i.Content,
		&i.MailingID,
		&i.InsertTime,
	)
	return i, err
}

const getByMailingID = `-- name: GetByMailingID :many
SELECT id, email, title, content, mailing_id, insert_time FROM messages 
where mailing_id = $1
`

func (q *Queries) GetByMailingID(ctx context.Context, mailingID int32) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getByMailingID, mailingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Title,
			&i.Content,
			&i.MailingID,
			&i.InsertTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}