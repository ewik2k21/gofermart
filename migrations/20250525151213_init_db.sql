-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       id uuid NOT NULL,
                       login text UNIQUE,
                       password_hash text,
                       PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
