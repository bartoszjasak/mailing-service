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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages
-- +goose StatementEnd
