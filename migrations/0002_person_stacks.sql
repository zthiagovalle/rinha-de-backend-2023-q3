-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS person_stacks (
    id BIGSERIAL PRIMARY KEY,
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    name VARCHAR(32) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS person_stacks;
-- +goose StatementEnd