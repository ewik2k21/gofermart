-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdraws (
                           user_id uuid NOT NULL PRIMARY KEY ,
                           "order" TEXT NOT NULL,
                           sum INTEGER DEFAULT 0,
                           processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdraws;
-- +goose StatementEnd
