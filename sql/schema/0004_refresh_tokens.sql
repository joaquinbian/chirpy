-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens(
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expired_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd


