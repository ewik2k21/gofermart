-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD CONSTRAINT unique_order_number UNIQUE (order_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP CONSTRAINT unique_order_number
-- +goose StatementEnd
