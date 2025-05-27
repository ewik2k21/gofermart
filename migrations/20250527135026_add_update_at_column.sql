-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN update_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

ALTER TABLE orders
    ALTER COLUMN id SET DEFAULT gen_random_uuid();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN updated_at;

ALTER TABLE orders
    ALTER COLUMN id DROP DEFAULT;
-- +goose StatementEnd
