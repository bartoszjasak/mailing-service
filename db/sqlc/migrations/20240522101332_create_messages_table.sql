-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
  id          SERIAL      PRIMARY KEY,
  email       text        NOT NULL,
  title       text        NOT NULL,
  content     text        NOT NULL,
  mailing_id  int         NOT NULL,
  insert_time timestamptz NOT NULL
);

CREATE INDEX messages_mailing_id_idx ON messages(mailing_id);
CREATE INDEX messages_insert_time_idx ON messages(insert_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX messages_insert_time_idx;
DROP INDEX messages_mailing_id_idx;
DROP TABLE messages;
-- +goose StatementEnd
