-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdraws DROP CONSTRAINT withdraws_pkey;
ALTER TABLE withdraws ADD COLUMN id UUID NOT NULL;

ALTER TABLE withdraws ADD PRIMARY KEY (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdraws
    DROP COLUMN id;
-- +goose StatementEnd
