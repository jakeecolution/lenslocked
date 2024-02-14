-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
