-- +goose Up
-- +goose StatementBegins
CREATE TABLE users (
                       id uuid PRIMARY KEY ,
                       login VARCHAR(255) NOT NULL UNIQUE ,
                       password_hash VARCHAR(255)
);



CREATE TABLE balances (
    user_id uuid PRIMARY KEY REFERENCES users(id),
    current FLOAT8 NOT NULL,
    withdraw FLOAT8 NOT NULL DEFAULT 0.00
);

CREATE TABLE orders (
    id uuid PRIMARY KEY ,
    user_id uuid NOT NULL REFERENCES users(id),
    order_number TEXT NOT NULL ,
    status VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
