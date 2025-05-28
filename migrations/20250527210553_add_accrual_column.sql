-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN accrual INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN accrual;
-- +goose StatementEnd
